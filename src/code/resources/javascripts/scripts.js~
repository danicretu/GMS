var tags=[];
var tagNo=0;
var tagTest=3;
var previousUrl="";
$(document).ready(function() {
	
	$('#menuButton').on('click', function(event){
						console.log("button pressed");
						$('div.sideMenu').removeClass('hidden-xs');
					});
	

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
	
	$("#deleteButton").click(function() {
		console.log("delete button was pressed");
		
		console.log("delete button was pressed");
		
		console.log("delete button was pressed");
		var name = $(this).attr("name");
		var picId =$("input#picId"+name).val();
		var albumId=$("input#picAlbumId"+name).val();
		var picOwner =$("input#picOwner"+name).val();
		var contentType = $("input#contentType"+name).val();
		console.log("delete button was pressed");
		$.ajax({
			url:"/delete",
			type:"POST",
			data:{"pic":picId, "album":albumId, "owner":picOwner, "cType":contentType},
			success: function(html){
				var t=html.split('_');
				if (t[0] == 'Yes') {
					
					console.log("in success "+html);
					console.log("in success "+t[1]);
					
						$('#li'+t[1]).remove();
						$('#picModal'+t[1]).remove();
				}
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
	
	$("#testButton").click(function(){
		console.log("in test button");
		$.ajax({
			url:"/flickrTest",
			type:"POST",
			success: function(html){
				console.log(html);
				document.getElementById('test').innerHTML = html;
				
				
			}
		});
		
	});
	
	$('.scroll-link').on('click', function(event){
		event.preventDefault();
		var sectionID = $(this).attr("data-id");
		scrollToID('#' + sectionID, 750);
	});	
	
	$("#commentForm").submit(function(){
		var comment =$("input#comment").val();
		var picture=$("input#pictureNumber").val();
		var album=$("input#albumNumber").val();
		var owner=$("input#owner").val();
		var cType=$("input#cType").val();
		$.ajax({
			type:"POST",
			url:"/saveComment",
			data:{"comment" : comment, "pic" : picture, "album":album, "owner":owner, "cType":cType},
			success: function(html) {
				var t=html.split('_');
				if (t[0]=='Yes') {
					$('#comment').val("");
					$('#commentList').prepend("<li>"+
									"<div class='commentText'>"+
									"<p>"+t[1]+"</p>"+
									"<a class='user under' href='#getUser?"+t[2]+"'>"+t[2]+"</a>"+
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
		assign(imgList[i].id)();
		
	}
	
	
	window.onpopstate = function(event) {
		console.log(location.pathname,"  ", location.hash);
		console.log(event.state);
		var func = location.hash.substring(1);
		var func_param = location.hash.substring(1).split("?");
		var myFunc = window[func_param[0]];
		if (func_param[1].indexOf(",") > -1){
			var args = func_param[1].split(',');
			console.log(args);
			myFunc.apply(this, args);
		} else { 
			myFunc(func_param[1]);
		}
	};


});


function onDelete(id, cType){
	console.log("delete button was pressed");
		console.log("delete button was pressed"+id);
		$.ajax({
			url:"/delete",
			type:"POST",
			data:{"pic":id, "cType":cType},
			success: function(html){
				console.log(html);
				var t=html.split('_');
				if (t[0] == 'Yes') {
					
					console.log("in success "+html);
					console.log("in success "+t[1]);
					
					$('#li'+t[1]).remove();
					$('#picModal'+t[1]).remove();
					tagCloud('User');
				}
			}
		});
	return false;
}


function getAlbumDetails(album,start,cType,nModP, nModN){
	//console.log(data);
	if (start != -1) {
		$.ajax({
			url:"/albums",
			type:"POST",
			data: {"albumId":album, "start":start, "cType":cType, "nModP":nModP, "nModN" : nModN},
			success: function(html){
				
				var obj = jQuery.parseJSON(html);
				//console.log(obj[0].Content);
				if (obj[0].Content != ""){
					$('#panelBodyContent').html(obj[0].Content);
					var uls = document.getElementsByName("lia")
					console.log(uls.length);
					if (uls.length >= 1) {
						
						jQuery.each(uls, function(index, value) {
						   
							assignClass(uls[0].id)();
					   });
					}
					
					carousel();
				}
			}
						
		});
	}
	return false;
	
}

function createNewAlbum(name){
	console.log("in create new album");
	n=$('input#'+name).val();
	$.ajax({
			url:"/createAlbum",
			type:"POST",
			data: {"name":n},
			success: function(html){
				$('#albumSelect')
		         .append($("<option></option>")
		         .attr("value",html)
		         .text(n)); 
				$('#albumModal').modal('hide');
				
			}
			});
	
}


function contentForm(){
		
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
					
					$('#contentType').val(document.querySelector('input[name="optradio"]:checked').value)
					if (t[3] != 'nil') {
						$('#imageLocation').attr('value', t[2]+","+t[3]);
						$('#find').click();
						$('#imageLocation').attr('value', $('#imageLocation').val());
					}
					console.log($('#contentType').val());
					if ($('#contentType').val() == "video"){
					document.getElementById('preview').style.display='none';
					}
					$("#uploadDiv").replaceWith($("#uploadForm"));
					document.getElementById('uploadForm').style.display = 'block';
					$("#uploadForm").keypress(function(e) {
						  //Enter key
						if (e.which == 13) {
						    return false;
						}
					});
					unbind();
				} else {
					console.log("fail upload");
				}
			}
		
		}); 
		return false;
}

function unbind(){
	$("#enterTag").unbind('keypress').keypress(function(e) {
			  //Enter key
			if (e.which == 13) {
				var t = $("input#enterTag").val();
				var tag = t.toLowerCase();
				console.log(tag,"   ",t);
				
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
}

function uploadForm(){
	var imageURL =$("input#imageURL").val();
	var caption=$("input#caption").val();
	var contentType=$("input#contentType").val();
	var e = document.getElementById("albumSelect");
	var albumSelect = e.options[e.selectedIndex].value;
	var location=$("input#location").val();
	var lng=$("input#lng").val();
	var lat=$("input#lat").val();
	var locality=$("input#locality").val();
	var formatted_address=$("input#formatted_address").val();
	var tags=$("input#tagList").val();
	console.log(imageURL);
	console.log(location+" location");
	console.log(lng+" lng");
	console.log(lat+" lat");
	console.log(locality+" locality");
	$.ajax({
		type:"POST",
		url:"/uploadPic",
		data:{"imageURL" : imageURL, "caption" : caption, "contentType":contentType, "albumSelect":albumSelect, "location":location, "lng":lng, "lat":lat, "locality":locality, "formatted_address" : formatted_address, "tagList":tags},
		success: function(html) {
			getPictures('');
			tagCloud('User');
		}
	}); 
	return false;
	
	
}

function getUpload(){
	setActive("uploadMenu");
	$.ajax({
		type:"GET",
		//url:"http://4e76fce3.ngrok.com/upload",
		url:"/upload",
		success: function(html) {

				$('div.sideMenu').removeClass('in');
				//console.log("in success"+html);
				var obj = jQuery.parseJSON(html);
				console.log(obj[0].Name+"            ***********")
				if (obj[0].Name=="upload") {
					console.log("in ok");
					console.log(document.getElementById('panelBodyContent').id);
					//document.getElementById('panelBodyContent').innerHtml = "Hello";
					$('#panelBodyContent').html(obj[0].Content);
					console.log(document.getElementById('panelBodyContent').innerHtml);
					
			}
				
		}
		
	});
	return false;
}

function commentFormSubmit(inp){
	console.log("in comment form");
	var comment =document.getElementById('comment'+inp).value;
	var picture=$("input#pictureNumber"+inp).val();
	var cType=$("input#cType"+inp).val();
	console.log("in comment form"+comment+" "+picture+" "+cType);
		$.ajax({
			type:"POST",
			url:"/saveComment",
			data:{"comment" : comment, "pic" : picture, "cType":cType},
			success: function(html) {
				var t=html.split('_');
				if (t[0]=='Yes') {
					$('#comment').val("");
					$('#commentList'+inp).prepend("<li>"+
									"<div class='commentText'>"+
									"<p>"+t[1]+"</p>"+
									"<a class='user under' href='#getUser?"+t[2]+"'>"+t[2]+"</a>"+
									"<span class='date under'> on "+t[3]+"</span>"+
									"</div></li>");
				} else {
					
				}
			}
		}); 
		return false;
}




function addTag(t, tagDiv) {
	var x = document.getElementById(tagDiv);
	for (var tag in t){
		var tagId = "tag"+tagNo++;
		console.log("in add tag t", t[tag])
		option = createElement(t[tag], tagDiv, tagId);

		x.appendChild(option);
	}
}

function createElement(tag, tagDiv, tagId){
	var option = document.createElement("a");
	option.text = tag;
		option.setAttribute('id',tagId);
		console.log("in addTag", tagId, " ", tagNo)
		if (tagDiv == "displayTags"){
			option.setAttribute('class', "tagUpload");
			option.setAttribute('onClick', function(event){removeTag(tagDiv);});
			option.onclick = function() {removeTag(tagDiv);};
			tags.push(tag);
			updateTagList();
		} else if (tagDiv == "suggestedTags") {
			option.setAttribute('class', "tagUpload");
			option.setAttribute('onClick', function(event){addToMainList(tag, tagNo);});
			option.onclick = function() {addToMainList(tag, tagNo);};
			
			
		}
		return option;
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

function addToMainList(tag, id) {

	var x = document.getElementById("displayTags");
	var option = document.createElement("a");
	var tagId = "tag"+id;
	removeTag("suggestedTags");
	console.log("in add to main list ", tag, id);
				
	option.text = tag;
	option.setAttribute('id',tagId);
	option.setAttribute('class', "tagUpload");
	option.setAttribute('onClick', function(event){removeTag("displayTags");});
	option.onclick = function() {removeTag("displayTags");};
	tags.push(tag);
	console.log("in add to main list", tags)
	updateTagList();
	x.appendChild(option);
}

function test(data){
	//var data =$("input#srch-term").val();
	console.log("in test"+data);
	return false;
	//document.getElementById('test').innerHTML = html;
}

function getVideos(data){
	console.log("in get Videos");
	setActive("videoMenu");
	if (data != -1){
		$.ajax({
			type:"POST",
			url:"/videos",
			data:{"req" : data},
			success : function(html){
				var obj = jQuery.parseJSON(html);
				console.log(obj[0].Name+"            ***********")
				if (obj[0].Name=="ownVideos") {
					$('div.sideMenu').removeClass('in');
					console.log("in ok");
					console.log(document.getElementById('panelBodyContent').id);
					//document.getElementById('panelBodyContent').innerHtml = "Hello";
					if (obj[0].Content != ""){
						$('#panelBodyContent').html(obj[0].Content);
						console.log(document.getElementById('panelBodyContent').innerHtml);
						
						var uls = document.getElementsByName("lia")
						console.log(uls.length);
						if (uls.length >= 1) {
							jQuery.each(uls, function(index, value) {
							    console.log(uls.length);
								assignClass(uls[0].id)();
					   		});
						}
						carousel();
					}
				}
				
			}
		})
	}
}

function getPictures(data){
	if (data!="-1"){
	console.log("in get Pictures");
	setActive("pictureMenu");
		$.ajax({
			type:"POST",
			url:"/pictures",
			data:{"req" : data},
			success: function(html) {
				var obj = jQuery.parseJSON(html);
				console.log(obj[0].Name+"            ***********")
				if (obj[0].Name=="ownPictures") {		
					$('div.sideMenu').removeClass('in');
					if (obj[0].Content != "") {		
					console.log("in ok");
					console.log(document.getElementById('panelBodyContent').id);
					//document.getElementById('panelBodyContent').innerHtml = "Hello";
					$('#panelBodyContent').html(obj[0].Content);
					var uls = document.getElementsByName("lia")
					console.log(uls.length);
					if (uls.length >= 1) {
						jQuery.each(uls, function(index, value) {
						    
							assignClass(uls[0].id)();
					      
					   	});
					}
					console.log(document.getElementById('panelBodyContent').innerHtml);
					carousel();
				}
			}
			}
	});
	return false;
	}
}

function assignClass(data) {

	return function() {

		document.getElementById(data).setAttribute("name","liaOwn");
		
	}

}

function assign(data) {
	console.log("in assign", data);
	return function() {
		document.getElementById(data).addEventListener("click",function(){ return upview(data);});
		
	}
}

function setClass(id) {
	$('#'+ data).attr('name','liaOwn');
}

function setActive(lid){
	var lis = document.getElementsByName("menuItem")
	jQuery.each(lis, function(index, value) {
		
		if (lis[index].id==lid) {
			$('#'+lis[index].id).addClass('active');
		}else{
			$('#'+lis[index].id).removeClass('active');
		}
	}); 
}

function getAlbums(data){
	console.log("in get Albums ", data);
	
	setActive("albumMenu");
	if (data == "") {
		$.ajax({
			type:"POST",
			url:"/albums",
			data:{"req" : data},
			success: function(html) {
				$('div.sideMenu').removeClass('in');
				console.log("in success"+html);
				var obj = jQuery.parseJSON(html);
				console.log(obj[0].Name+"            ***********")
				if (obj[0].Name=="ownAlbums") {
					console.log("in ok");
					console.log(document.getElementById('panelBodyContent').id);
					//document.getElementById('panelBodyContent').innerHtml = "Hello";
					$('#panelBodyContent').html(obj[0].Content);
					console.log(document.getElementById('panelBodyContent').innerHtml);
				}
				
			}, error: function(data) {
                console.log(data);
            }
	});
	return false;
	}
}

function flickrMenu(data, start, cType){
	if (data=="getTags"){
		data +=$("input#srch-term").val();
	} 
	var init = $("input#srch-term").val();
	setActive("flickrMenu");
	if (start != -1){
		$.ajax({
			type:"POST",
			url:"/flickrImages",
			data:{"req":data, "init":init, "start":start, "cType":cType},
			success: function(html){
				if (data=="start"){
					document.getElementById('panelBodyContent').innerHTML=html;
					document.getElementById('resultOr').style.visibility='hidden';
				} else if(data.indexOf("getTags") > -1){
							document.getElementById('resultOr').style.visibility='hidden';
							console.log("in else")
							if (data.indexOf("getTags_") > -1){
								console.log(data.substring(8));
								$('input#srch-term').val(data.substring(8));
							}
							if (html != "No content found with requested tag"){
								var obj = jQuery.parseJSON(html);
								document.getElementById('cloudFlickr').innerHTML="";
								populateCloud(obj[0].Content,"Flickr2");
								document.getElementById('cloudFlickr').style.visibility='visible';
								document.getElementById('cloudHeader').style.visibility='visible';
								document.getElementById('resultAnd').innerHTML=obj[1].Content;
							} else {
								document.getElementById('cloudFlickr').innerHTML=html;
								document.getElementById('cloudFlickr').style.visibility='visible';
								document.getElementById('resultAnd').style.visibility='hidden';
							}
				} else {
					var obj = jQuery.parseJSON(html);
					
					if (cType == "and" && obj[0].Content!=""){
						document.getElementById('resultAnd').innerHTML=obj[0].Content;
						document.getElementById('resultAnd').style.visibility='visible';
					} else if (cType=="or" && obj[1].Content!="") {
						document.getElementById('resultOr').innerHTML=obj[1].Content;
						document.getElementById('resultOr').style.visibility='visible';
					}else{
						document.getElementById('resultAnd').innerHTML=obj[0].Content;
						document.getElementById('resultAnd').style.visibility='visible';
						document.getElementById('resultOr').innerHTML=obj[1].Content;
						document.getElementById('resultOr').style.visibility='visible';
					} 
				}
			carousel();
			}
		});
	}
	return false;
}

function flickrNews(data, start, cType){
	if (data=="getTags"){
		data +=$("input#srch-term").val();
	}
	console.log(data+"start");
	console.log(cType);
	setActive("flickrCWG");
	if (start != -1){
		$.ajax({
				type:"POST",
				url:"/flickrCwg",
				data:{"req" : data, "start" : start, "cType" : cType},
				success: function(html) {
					$('div.sideMenu').removeClass('in');
					console.log(html);
					if (data=="start"){
						var obj = jQuery.parseJSON(html);
						document.getElementById('panelBodyContent').innerHTML=obj[0].Content;
						document.getElementById('cloudFlickr').innerHTML="";
						populateCloud(obj[1].Content,"Flickr");
						document.getElementById('cloudFlickr').style.visibility='visible';
					} else {
						var obj = jQuery.parseJSON(html);
						console.log(obj+"            ***********")
						
							if (cType == "image" && obj[0].Content!=""){
								document.getElementById('resultPhotos').innerHTML=obj[0].Content;
								document.getElementById('resultPhotos').style.visibility='visible';
							} else if (cType=="news" && obj[1].Content!="") {
								document.getElementById('resultNews').innerHTML=obj[1].Content;
								document.getElementById('resultNews').style.visibility='visible';
							}else{
								document.getElementById('resultPhotos').innerHTML=obj[0].Content;
								document.getElementById('resultPhotos').style.visibility='visible';
								document.getElementById('resultNews').innerHTML=obj[1].Content;
								document.getElementById('resultNews').style.visibility='visible';
							} 
						
							
						
					}
				}
		});
	}
	return false;
}


function updateTagList() {
	var tagsForHTML = document.getElementById("tagList");
	console.log("in update tag list", tags);
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
    
function tagCloud(cloud) {
	
	
	$.ajax({
        url: "/tagCloud",
        type: "GET",
        success: function (data) {
			console.log(data);
			populateCloud(data, cloud);			
        },
            error: function(data) {
                console.log("Error getting tags from db");
            }
    });	
}

function getSimilarTag(t,start,cType,nModP, nModN){
	console.log(t);
	setActive("tags");
	$.ajax({
        url: "/tag",
        type: "GET",
		data: {"tag" : t,"start":start, "cType":cType, "nModP":nModP, "nModN":nModN},
        success: function (html) {
				if (html != ""){
					$('div.sideMenu').removeClass('in');
					$('#panelBodyContent').html(html);
					//document.getElementById('deleteButton').style.display='none';
					var imgList = document.getElementsByName("lia");
					carousel();
					for (var i = 0; i < imgList.length; i++)
					{
						assign(imgList[i].id)();
						
					}
				}
				
        		},
            error: function(data) {
                console.log("Error getting tags from db");
            }
    });
}

function populateCloud(data, cloud){
	var tagMap = {};
	$('#cloud'+cloud).html("");
	console.log("cloud "+data)
	var t=data.split(',');
	
	if (cloud=="Flickr2"){
		for (var m in t){
			$('#cloudFlickr').append("<a class='tag' href='#flickrMenu?" + t[m] + "'>"+t[m]+"</a>");
		}
		return;
	}
	
	var max = parseInt(t.pop().split(' ')[1]);
	for (i=0; i<t.length; i++) {
		var split=t[i].split(' ');
		tagMap[split[0]]=parseInt(split[1]);
	}
	console.log(tagMap+"      "+max);
	
	for (var m in tagMap){
		if(tagMap[m] > 0){
			
			if(tagMap[m]/max == 1) size = 8;
			else if((1>tagMap[m]/max) && (tagMap[m]/max>0.7)) size = 7;
			else if((0.7>tagMap[m]/max) && (tagMap[m]/max>0.5)) size = 6;
			else if ((0.5>tagMap[m]/max) && (tagMap[m]/max>0.3)) size = 4;
			else size = 2;
			if (cloud=="Flickr"){
				$('#cloud'+cloud).append("<a class='size"+size+"' href='#flickrNews?" + m + "'>"+m+"</a>");
			}else{
				$('#cloud'+cloud).append("<a class='size"+size+"' href='#getSimilarTag?" + m + "'>"+m+"</a>");
			}
		}
	}
}

function getUser(u,start,cType,nModP, nModN){
	setActive("users");
	$.ajax({
        url: "/user",
        type: "GET",
		data:{"user":u, "start":start, "cType":cType, "nModP":nModP, "nModN":nModN},
        success: function (data) {
			console.log(data);
			if (data!=""){
				$('#panelBodyContent').html(data);
				var imgList = document.getElementsByName("lia");
				for (var i = 0; i < imgList.length; i++)
				{
					assign(imgList[i].id)();
					
				}
				//document.getElementById('deleteButton').style.display='none';
				carousel();
			}	
        },
        error: function(data) {
                console.log("Error getting tags from db");
        }
    });	
}

function scrollToID(id, speed){
	var offSet = 50;
	var targetOffset = $(id).offset().top - offSet;
	var mainNav = $('#main-nav');
	$('html,body').animate({scrollTop:targetOffset}, speed);
	if (mainNav.hasClass("open")) {
		mainNav.css("height", "1px").removeClass("in").addClass("collapse");
		mainNav.removeClass("open");
	}
}

function checkIfLoggedIn() {
	console.log("in check login");
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
					var uls = document.getElementsByName('logP');
					for (var i = 0; i < uls.length; i++){
						
						uls[i].style.display='none';
					}
					var uls1 = document.getElementsByName('logR');
					for (var i = 0; i < uls1.length; i++){
						
						uls1[i].style.display='none';
					}
					
					var retrieveUsers = document.getElementsByName('logA');
					for (var i=0; i<retrieveUsers.length; i++){
						 
						retrieveUsers[i].href="#";
					}
					
					var tags = document.getElementsByName('logT');
					for (var i=0; i<tags.length; i++){
						 
						tags[i].href="#";
					}
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
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+next.id);
				//document.getElementById("next"+ul[m].id).addEventListener("click",function(){ print(next.id); });
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+prev.id);
				//document.getElementById("prev"+ul[m].id).addEventListener("click",function(){print(prev.id); });
				addListener(ul[m].id, next.id, prev.id)();
				
			} else if (m==0){
				
				var next1 = ul[m+1]
				var prev1 = ul[ul.length-1]
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+next1.id);
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+prev1.id);
				//document.getElementById("next"+ul[m].id).addEventListener("click",function(){ print(next1.id); });
				//document.getElementById("prev"+ul[m].id).addEventListener("click",function(){ print(prev1.id); });
				addListener(ul[m].id, next1.id, prev1.id)();
			} else {
				
				var next2 = ul[m+1]
				var prev2 = ul[m-1]
				document.getElementById("next"+ul[m].id).setAttribute('data-target','#picModal'+next2.id);
				document.getElementById("prev"+ul[m].id).setAttribute('data-target','#picModal'+prev2.id);
				//document.getElementById("next"+ul[m].id).addEventListener("click",function(){ print(next2.id); });
				//document.getElementById("prev"+ul[m].id).addEventListener("click",function(){ print(prev2.id); });
				addListener(ul[m].id, next2.id, prev2.id)();
			}
		}
	}
	
	var ul2 = document.getElementsByName("liaOwn");
	if (ul2.length >1){
		for (m=0; m<ul2.length; m++) {
			if (m==ul2.length-1){
				var next3 = ul2[0]
				var prev3 = ul2[m-1]
				document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+next3.id);
				document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+prev3.id);
			} else if (m==0){
				var next4 = ul2[m+1]
				var prev4 = ul2[ul2.length-1]
				document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+next4.id);
				document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+prev4.id);
			} else {
				var next5 = ul2[m+1]
				var prev5 = ul2[m-1]
				document.getElementById("next"+ul2[m].id).setAttribute('data-target','#picModal'+next5.id);
				document.getElementById("prev"+ul2[m].id).setAttribute('data-target','#picModal'+prev5.id);
			}
		}
	}
}

function addOnClick(aLink, tag, cloud){
	
	if (cloud=="Flickr"){
		return function() {
			aLink.onclick=function() {return flickrNews(tag);};
		}
	}else{
		return function() {
			aLink.onclick=function() {return flickrMenu(tag);};
		}
	}
}


function addListener(element, next, prev){
	return function() {
		document.getElementById("next"+element).addEventListener("click",function(){ upview(next); });
		document.getElementById("prev"+element).addEventListener("click",function(){ upview(prev); });
	}
}

function upview(id) {
		var picId =$("input#picId"+id).val();
		var albumId=$("input#picAlbumId"+id).val();
		var picOwner =$("input#picOwner"+id).val();
		var cType=$("input#contentType"+id).val();
		$.ajax({
			url:"/upvote",
			type:"POST",
			data:{"picId" : picId, "albumId" : albumId, "picOwner" : picOwner, "cType":cType},
			success: function(html){
				var t=html.split('_');
				if (t[0] == 'Yes') {
					document.getElementById("upvotep"+id).innerHTML = t[1];
				}
				
			}
	}); 


}
