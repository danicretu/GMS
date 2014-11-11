$(document).ready(function() {
	
	$("#loginForm").submit(function() {
		//document.getElementById('errorLabel').innerHTML('no error');
		var email =$("input#email_label").val();
		var password=$("input#password_label").val();
		$.ajax({
			type:"POST",
			url:"/login",
			data:{"email" : email, "pass" : password},
			success: function(html) {
				if (html=='Yes') {
					console.log("**********************")
					setTimeout('go_to_userPage()', 500);
				} else {
					document.getElementById("errorLabel").innerHTML="wrong username or password";
				}
			}
		});
		
		return false;
	});
	
	
});

function go_to_userPage() {
	window.location="/auth"
}


/*
$(document).ready(function(){
	
  $(".error").hide();
  $("#loginForm").submit(function(){
	var email = $("input#email_label").val();
	var password=$("input#password_label").val();
    $.ajax({
			type:"POST"
			url:"/login",
			data:{"email":email, "pass":password}, 
			success: function(html) {
					if (html=='Yes') {
						console.log("****************************");
						setTimeout('go_to_userPage()', 3000);
						
					} else {
						document.getElementById('errorLabel').innerHTML='wrong username or password';
					
					}
				}});
	return false;
  }
);

	
}); 

function go_to_userPage() {
	window.location="/auth"
}

*/

/*$(function() {
	$(".error").hide();
	
	$(".userLoginButton").click(function() {
	
	//	var email = $("input#email_label").val();
	//	var password=$("input#password_label").val();
	//	var data = email+" "+password;
		
		$ajax({
			type:"POST",
			url:"http://localhost:8080/login",
	//		data: data,
			success:function() {
				$("#wrongCredentials").html("<p>Wrong username or password</p>");
			}
			
		});
		return false;
	
	});
	
	
}); */