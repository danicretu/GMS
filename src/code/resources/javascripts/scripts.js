var tags=[];
var tagNo=0;
var tagTest=3;
$(document).ready(function() {

	$("#imgInp").change(function(){
		//console.log("=====================================");
	    //readURL(this);
		//document.getElementById('submit').style.visibility = 'visible';
	});
	
	$("#albumSubmit").click(function(){
		var name =$("input#albumName").val();
		var description=$("input#albumDescription").val();
		$.ajax({
			url:"/createAlbum",
			type:"POST",
			data:{"name" : name, "description" : description},
			success: function(html){
				$('#albumSelect').append("<option value="+html+">"+name+"</option>");
				document.getElementById('albumModal').style.visibility = 'hidden';
			}
		});
		
	});
	
	$("#upvote").click(function(){
		var picId =$("input#picId").val();
		var albumId=$("input#picAlbumId").val();
		var picOwner =$("input#picOwner").val();
		$.ajax({
			url:"/upvote",
			type:"POST",
			data:{"picId" : picId, "albumId" : albumId, "picOwner" : picOwner},
			success: function(html){
				var t=html.split('_');
				if (t[0] == 'Yes') {
					document.getElementById("upvotep").innerHTML = t[1];
				}
				
			}
		}); 
		
	});
	
	$("#submit").click(function(){
		
		file=$('#imgInp')[0].files[0];
		var form = new FormData();
		form.append("uploadData", file);
		
		$.ajax({
			url:"/saveFile",
			type:"POST",
			data: form,
			processData: false,
			contentType: false,
			success: function(html){
				var t=html.split('_');
				if (t[0]=='Yes') {
					$('#blah').attr('src', t[1]);
					$('#imageURL').val(t[1]);
					if (t[3] != 'nil') {
						$('#imageLocation').attr('value', t[2]+","+t[3]);
						$('#find').click();
						$('#imageLocation').attr('value', $('#imageLocation').val());
					}
					$("#uploadDiv").replaceWith($("#uploadForm"));
					document.getElementById('uploadForm').style.display = 'block';
				} else {
					console.log("fail upload");
				}
			}
		
		});
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
			
			if (tag != "" && tag != " ") {
				if (document.getElementById('tagsLabel').style.visibility == 'hidden'){
					document.getElementById('tagsLabel').style.visibility = 'visible';
				}
				var t=tag.split(',');
				addTag(t, "displayTags");
				if (document.getElementById('displayTags').style.visibility == 'hidden'){
					document.getElementById('displayTags').style.visibility = 'visible';
				}

				flickrRelatedTags(tag);
				//tagAlgo(tag)
				
			}
			$('#enterTag').val("");
		}
	});
	
	$("#commentForm").submit(function(){
		var comment =$("input#comment").val();
		var picture=$("input#pictureNumber").val();
		var album=$("input#albumNumber").val();
		var owner=$("input#owner").val();
		$.ajax({
			type:"POST",
			url:"/saveComment",
			data:{"comment" : comment, "pic" : picture, "album":album, "owner":owner},
			success: function(html) {
				var t=html.split('_');
				if (t[0]=='Yes') {
					$('#comment').val("");
					$('#commentList').prepend("<li>"+
									"<div class='commentText'>"+
									"<p>"+t[1]+"</p>"+
									"<a class='user under' href='/user?"+t[2]+"'>"+t[2]+"</a>"+
									"<span class='date under'> on "+t[3]+"</span>"+
									"</div></li>");
				} else {
					
				}
			}
		});
		return false;
	});
	
	var imgList = document.getElementsByName("lia");
	for (var i = 0; i < imgList.length; i++)
	{
		assign(imgList[i]);
		//var idl = imgList[i]
		//console.log(idl.id);
	    //imgList[i].addEventListener("click",function(data){ test(data); },idl.id);
	}

});

function assign(data) {
	data.addEventListener("click",function(){ test(data.id);});
}


function addTag(t, tagDiv) {
	var x = document.getElementById(tagDiv);
	for (var tag in t){
		var option = document.createElement("a");
		var tagId = "tag"+tagNo++;
					
		option.text = t[tag];
		option.setAttribute('id',tagId);
		
		if (tagDiv == "displayTags"){
			option.setAttribute('class', "tagUpload");
			option.setAttribute('onClick', function(event){removeTag(tagDiv);});
			option.onclick = function() {removeTag(tagDiv);};
			tags.push(t[tag]);
			updateTagList();
		} else if (tagDiv == "suggestedTags") {
			option.setAttribute('class', "tagUpload");
			option.setAttribute('onClick', function(event){addToMainList(t[tag]);});
			option.onclick = function() {addToMainList(t[tag]);};
			
			
		}
		x.appendChild(option);
	}
}

function removeTag(list) {
	var text = $(event.target).text();
	var index = jQuery.inArray(text,tags);
	var tagList = document.getElementById(list);
	var tag = document.getElementById(event.target.id);
	tagList.removeChild(tag);
	if (index != -1) {
		tags.splice(index, 1);
		updateTagList();
	}
}

function addToMainList(tag) {
	var x = document.getElementById("displayTags");
	var option = document.createElement("a");
	var tagId = "tag"+tagNo++;
	removeTag("suggestedTags");
				
	option.text = tag;
	option.setAttribute('id',tagId);
	option.setAttribute('class', "tagUpload");
	option.setAttribute('onClick', function(event){removeTag("displayTags");});
	option.onclick = function() {removeTag("displayTags");};
	tags.push(tag);
	updateTagList();
	x.appendChild(option);
}

function updateTagList() {
	var tagsForHTML = document.getElementById("tagList");
	tagsForHTML.setAttribute('value', tags);
}



function flickrRelatedTags(tag) {
	var url1 = "https://api.flickr.com/services/rest/?method=flickr.tags.getRelated&api_key=ef72e911f885e924a460b98a4801ff14&tag=";
	var url2 = "&per_page=5&format=json";
	$.ajax({
        url: "/flickr",
        type: "GET",
		data: {"url1":url1,"url2":url2, "tags":tag},
        success: function (data) {
			var myNode = document.getElementById("suggestedTags");
			while (myNode.firstChild) {
   				 myNode.removeChild(myNode.firstChild);
			}

            processFlickrTags(data)
        },
            error: function(data) {
                var err = ("(" + xhr.responseText + ")");
            }
    });
}

function processFlickrTags(tags) {
	var indivTags = tags.split(',');
	indivTags.pop();

	indivTags = indivTags.slice(0,10)
	addTag(indivTags, "suggestedTags")

	if (document.getElementById('suggestedTags').style.visibility == 'hidden'){
		document.getElementById('suggestedTags').style.visibility = 'visible';
	}
	
	document.getElementById('tagsLabel').style.visibility = 'hidden';
}




function readURL(input) {

        if (input.files && input.files[0]) {
            var reader = new FileReader();
            
            reader.onload = function (e) {
                $('#blah').attr('src', e.target.result);
				document.getElementById('blah').style.visibility='visible';
				document.getElementById('photoDetails').style.visibility='visible';
            }
            
            reader.readAsDataURL(input.files[0]);
        }
    }
    
function tagCloud() {
	var tagMap = {};
	$.ajax({
        url: "/tagCloud",
        type: "GET",
        success: function (data) {

			var t=data.split(',');
			var max = parseInt(t.pop().split(' ')[1]);
			for (i=0; i<t.length; i++) {
				var split=t[i].split(' ');
				tagMap[split[0]]=parseInt(split[1]);
			}
			for (var m in tagMap){
				if(tagMap[m] > 0){

					if(tagMap[m]/max == 1) size = 8;
					else if((1>tagMap[m]/max) && (tagMap[m]/max>0.7)) size = 7;
					else if((0.7>tagMap[m]/max) && (tagMap[m]/max>0.5)) size = 6;
					else if ((0.5>tagMap[m]/max) && (tagMap[m]/max>0.3)) size = 4;
					else size = 2;
					$('#cloud').append("<a class='size"+size+"' href='/tag?"+m+"'>"+m+"</a>");
				}
			}
			
        },
            error: function(data) {
                console.log("Error getting tags from db");
            }
    });	
}

function checkIfLoggedIn() {
	$.ajax({
			type:"GET",
			url:"/checkLogIn",
			success: function(html) {
				var t=html.split(',');
				if (t[0]=='Yes') {
					$('#loggedIn').attr('class', 'dropdown');
					document.getElementById('loggedIn').innerHTML='<a href="/authenticated2" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">'+t[1]+'<span class="caret"></span></a>'+
																'<ul class="dropdown-menu" role="menu">'+
																	'<li><a href="/authenticated">Profile</a></li>'+
																	'<li><a href="/logout">Log Out</a></li></ul>';
				} else {
					document.getElementById('loggedIn').innerHTML='<a href="#" data-toggle="modal" data-target="#loginModal">Log In</a>';
				}
			}
		});
}

function carousel() {
	var ul = document.getElementsByName("lia");
	if (ul.length > 1) {
		for (m=0; m<ul.length; m++) {

			if (m==ul.length-1){
				var next = ul[0]
				var prev = ul[m-1]
				console.log(next.id+" 1next")
				console.log(prev.id+" 1prev")
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+ul[0].id);
				document.getElementById("next"+ul[m].id).addEventListener("click",function(){ test(next.id); });
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+ul[m-1].id);
				document.getElementById("prev"+ul[m].id).addEventListener("click",function(){ test(prev.id); });
			} else if (m==0){
				var next1 = ul[m+1]
				var prev1 = ul[ul.length-1]
				console.log(next1.id+" 2next")
				console.log(prev1.id+" 2prev")
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+ul[m+1].id);
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+ul[ul.length-1].id);
				document.getElementById("next"+ul[m].id).addEventListener("click",function(){ test(next1.id); });
				document.getElementById("prev"+ul[m].id).addEventListener("click",function(){ test(prev1.id); });
			} else {
				var next2 = ul[m+1]
				var prev2 = ul[m-1]
				console.log(next2.id+" 3next")
				console.log(prev2.id+" 3prev")
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+ul[m+1].id);
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+ul[m-1].id);
				document.getElementById("next"+ul[m].id).addEventListener("click",function(){ test(next2.id); });
				document.getElementById("prev"+ul[m].id).addEventListener("click",function(){ test(prev2.id); });
			}
		}
	}
	
	var ul2 = document.getElementsByName("liaOwn");
	for (m=0; m<ul2.length; m++) {
		if (m==ul2.length-1){
			document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[0].id);
			document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[m-1].id);
		} else if (m==0){
			document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[m+1].id);
			document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[ul2.length-1].id);
		} else {
			document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[m+1].id);
			document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+ul2[m-1].id);
		}
	}
}

function test(id) {
		console.log(id)
		var picId =$("input#picId"+id).val();
		var albumId=$("input#picAlbumId"+id).val();
		var picOwner =$("input#picOwner"+id).val();
		$.ajax({
			url:"/upvote",
			type:"POST",
			data:{"picId" : picId, "albumId" : albumId, "picOwner" : picOwner},
			success: function(html){
				var t=html.split('_');
				if (t[0] == 'Yes') {
					document.getElementById("upvotep"+id).innerHTML = t[1];
				}
				
			}
	}); 
	console.log("in test function");

}
