var tags=[];
var tagNo=0;
$(document).ready(function() {
	
	
	
	$("#imgInp").change(function(){
		console.log("=====================================");
	    readURL(this);
		document.getElementById('submit').style.visibility = 'visible';
	});
	
	$("#uploadForm").keypress(function(e) {
		  //Enter key
		if (e.which == 13) {
		    return false;
		}
	});
	
	$("#enterTag").unbind('keypress').keypress(function(e) {
		  //Enter key
		if (e.which == 13) {
			var tag = $("input#enterTag").val();
			console.log("entering tag");
			console.log(tag);
			
			if (tag != "" && tag != " ") {
				var x = document.getElementById("displayTags");
		    		var option = document.createElement("a");
				var tagId = "tag"+tagNo++;
				
		    		option.text = tag;
				option.setAttribute('id',tagId);
				option.setAttribute('class', "tagUpload");
				option.setAttribute('onClick', function(event){removeTag();});
				option.onclick = function() {removeTag();};
				tags.push(tag);
		    		x.appendChild(option);
				console.log(tags);
				if (document.getElementById('displayTags').style.visibility == 'hidden'){
					document.getElementById('displayTags').style.visibility = 'visible';
				}
				
				
				
			}
			$('#enterTag').val("");
		}
	});
	
/*	$("#tag").click(function(event) {
		
		var text = $(event.target).text();
		if (jQuery.inArray(text, tags)) {
			
		}
	
	}); */
	
});

function removeTag() {
	console.log("working");
	var text = $(event.target).text();
	var index = jQuery.inArray(text,tags);
	if (index != -1) {
		console.log(event.target.id);	
		console.log(index);
		var tagList = document.getElementById("displayTags");
		var tag = document.getElementById(event.target.id);
		tagList.removeChild(tag);
		tags.splice(index, 1);
		console.log(tags);	
	}
}



function readURL(input) {
		console.log("///////////////////////////////////////")
		console.log(input);
        if (input.files && input.files[0]) {
            var reader = new FileReader();
            
            reader.onload = function (e) {
                $('#blah').attr('src', e.target.result);
				document.getElementById('blah').style.visibility='visible';
				document.getElementById('photoDetails').style.visibility='visible';
				console.log(e.target.result);
            }
            
            reader.readAsDataURL(input.files[0]);
        }
    }
    
