$(document).ready(function() {
	
	$("#login").click(function() {
		console.log("in login");
		el = document.getElementById("loginBox");
		el.style.visibility = (el.style.visibility == "visible") ? "hidden" : "visible";
	});
	
	$("#close").click(function() {
		console.log("in close");
		el = document.getElementById("loginBox");
		el.style.visibility = (el.style.visibility == "visible") ? "hidden" : "visible";
	});
	
	$("#loginForm").submit(function() {
		//document.getElementById('errorLabel').innerHTML('no error');
		var email =$("input#email_label").val();
		var password=$("input#password_label").val();
		console.log("in login");
		$.ajax({
			type:"POST",
			url:"/login",
			data:{"email" : email, "pass" : password},
			success: function(html) {
				console.log("**********************");
				var t=html.split("_");
				
				if (t[0]=='Yes') {
					console.log("**********************");
					//setTimeout('go_to_userPage()', 500);
					go_to_userPage();
					
		            
				} else {
					document.getElementById("errorLabel").innerHTML="wrong username or password";
				}
			}
		});
		
		return false;
	});
	
	
	$("#registerButton").click(function(){
		$("#modalBody").replaceWith($("#registerDiv"));
		document.getElementById('registerDiv').style.display='block';
		document.getElementById('modalLabel').innerHTML='Register for an Account';
	});
	
	$("#registerForm").submit(function() {
		//document.getElementById('errorLabel').innerHTML('no error');
		var email =$("input#email").val();
		var password=$("input#password").val();
		var password2=$("input#confirmPassword").val();
		var last=$("input#last").val();
		var first=$("input#first").val();
		
		if (password != password2) {
			document.getElementById("errorLabel").innerHTML="Passwords not identical";
		} else {
		
			$.ajax({
				type:"POST",
				url:"/register",
				data:{"email" : email, "pass" : password, "first" : first, "last" : last},
				success: function(html) {
					if (html=='Yes') {
						console.log("**********************")
						setTimeout('go_to_userPage()', 500);
					}
				}
			});
			
		}
		
		return false;
	});
	

	
});

function go_to_userPage() {
	window.location="/authenticated";
}

function registerBox() {
	el = document.getElementById("registerBox");
	el.style.visibility = (el.style.visibility == "visible") ? "hidden" : "visible";
}

