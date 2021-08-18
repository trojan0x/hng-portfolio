var i = 0;
	var txt = "My name is Jeff Dauda, a tech guy... I live in Jos, Plateau state in Nigeria"; /* The text */
	var speed = 50; /* The speed/duration of the effect in milliseconds */
	
	function typeWriter() {
	  if (i < txt.length) {
		document.getElementById("jeff").innerHTML += txt.charAt(i);
		i++;
		setTimeout(typeWriter, speed);
	  }
	}