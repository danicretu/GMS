import java.sql.DriverManager;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.TreeMap;


public class Recommendation {
	
	public static void main(String args[]){
		Connection connect = null;
		
		String url="jdbc:mysql://localhost:3306/testDB";
		String user="root";
		String password="";
		String topFiveTags="";
		String startTag="";
		Map<String, Double> frequencyMap = new HashMap<String, Double>();
		
		if (args.length == 1) {
			startTag=args[0].trim();
			System.out.println("input tag: "+startTag);		
			try{
				connect = DriverManager.getConnection(url, user, password);
				Statement st= connect.createStatement();
				ResultSet rs = st.executeQuery("SELECT * FROM cooc_vectors WHERE tag='"+startTag+"'");
				
			
				if (rs.next()) {
					String frequencies = rs.getString(2);
					String[] tagList = frequencies.split(" ");
					for (int i = 0; i<tagList.length; i++) {
						String[] tags = tagList[i].split(":");
						frequencyMap.put(tags[0], Double.parseDouble(tags[1]));
						
					}
					
					for (String tag : frequencyMap.keySet()) {
						rs = st.executeQuery("Select idfimages from tags where pos='"+tag+"'");
						if (rs.next()) {
							double idf = Double.parseDouble(rs.getString(1));
							frequencyMap.put(tag, frequencyMap.get(tag)*idf);
						}
					}
	
					List<Map.Entry<String, Double>> sorted = sort(frequencyMap);
					for (int i = 0; i < 5; i++) {
						String tag = sorted.get(sorted.size()-i-1).getKey();
						rs = st.executeQuery("Select tag from tags where pos='"+tag+"'");
						if (rs.next()) {
							topFiveTags+=rs.getString(1)+" ";
						}
					}
					
					System.out.println(topFiveTags);
					
				}
				
				connect.close();
			} catch (SQLException e){
				System.out.println("SQL Exception ");
				e.printStackTrace();
			}
		} else if (args.length > 1) {
			for (int i = 0; i<args.length; i++) {
				System.out.println(args[i]);
			}

		}
	}
	
	private static List<Map.Entry<String, Double>> sort(Map<String, Double> map) {
		Comparator<Map.Entry<String, Double>> compareValues = new Comparator<Map.Entry<String, Double>>(){

			@Override
			public int compare(Map.Entry<String, Double> first,
					Map.Entry<String, Double> second) {
				return first.getValue().compareTo(second.getValue());
			}
			
		};
		
		List<Map.Entry<String, Double>> tagList = new ArrayList<Map.Entry<String, Double>>(map.entrySet());
		Collections.sort(tagList, compareValues);
		return tagList;
	}

	
}
