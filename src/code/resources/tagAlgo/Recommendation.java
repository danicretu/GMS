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
				try {
				connect = DriverManager.getConnection(url, user, password);
				System.out.println(processMultipleTags(args, connect));
				
				
				
			} catch (Exception e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
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
	
	private static Map<String, Double> populateTagMap(Connection connect, String startTag) throws Exception{
		Statement st= connect.createStatement();
		ResultSet rs = st.executeQuery("SELECT * FROM cooc_vectors WHERE tag='"+startTag+"'");
		Map<String, Double> freqMap = new HashMap<String, Double>();
	
		if (rs.next()) {
			String frequencies = rs.getString(2);
			String[] tagList = frequencies.split(" ");
			for (int i = 0; i<tagList.length; i++) {
				String[] tags = tagList[i].split(":");
				freqMap.put(tags[0], Double.parseDouble(tags[1]));
				
			}
		}
		return freqMap;
	}
	
	private static String processMultipleTags(String[] args, Connection connect) throws Exception {
		Statement st = connect.createStatement();
		Map<String, Double> allHit = new HashMap<String, Double>();
		Map<String, Double> others = new HashMap<String, Double>();
		Map<String, Double> main = new HashMap<String, Double>();

		List<Map<String, Double>> tagList = new ArrayList<Map<String, Double>>();
		for (int i=0; i< args.length; i++) {
			tagList.add(populateTagMap(connect, args[i]));
		}

		
		
		for (Map<String, Double> map : tagList){
			for (String key : map.keySet()) {
				main.put(key, 0.0);
			}
		}
		
		System.out.println("main size: "+main.size());
		
		for (String key : main.keySet()) {
			boolean isAll = true;
			for (Map<String, Double> map:tagList){
				if (map.keySet().contains(key)) {
					main.put(key, main.get(key) + map.get(key));
				} else {
					isAll = false;
				}
					
			}
			if (isAll == false) {
				others.put(key, main.get(key));
			} else {
				allHit.put(key, main.get(key));
			}
		}
		
		if (allHit.size() > 5) {
			System.out.println("allHit size: "+allHit.size());
			ResultSet rs = null;
			String topFiveTags = "";
			System.out.println(allHit);
			allHit = multiplyIDF(allHit, st);
			System.out.println(allHit);
			List<Map.Entry<String, Double>> sorted = sort(allHit);
			for (int i = 0; i < 5; i++) {
				String tag = sorted.get(sorted.size()-i-1).getKey();
				System.out.println(tag);
				rs = st.executeQuery("Select tag from tags where pos='"+tag+"'");
				if (rs.next()) {
					topFiveTags+=rs.getString(1)+" ";
				}
			}
			return topFiveTags;
			
		} else {
			allHit = multiplyIDF(allHit, st);
			others = multiplyIDF(others, st);
		}
		
		
		return "";
	}
	
	private static Map<String, Double> multiplyIDF(Map<String, Double> map, Statement st) throws SQLException {
		ResultSet rs=null;
		for (String tag : map.keySet()) {
			rs = st.executeQuery("Select idfimages from tags where pos='"+tag+"'");
			if (rs.next()) {
				double idf = Double.parseDouble(rs.getString(1));
				map.put(tag, map.get(tag)*idf);
			}
		}
		
		return map;
	}

	
}
