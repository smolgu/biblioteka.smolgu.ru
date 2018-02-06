var belorusskie;

$.getJSON("/static/js/bel-data.js?_=" + new Date().getTime(), function(data){
	belorusskie = data
	//$.getJSON("/bel/ssid", function(resp){
	//	ruslan_session = resp.session_id;
		init_bel();
	//})
})


function init_bel() {
	bel_draw_letter("а")
	$(".bel .letters a").click(function(e){
		$(".bel .results").html("")
		var letter = e.target.text.toLowerCase()
		bel_draw_letter(letter)
		e.preventDefault();
		return false;
	})
}

function bel_draw_letter(letter) {
	for (var i = 0; i < belorusskie.length; i++) {
			if (belorusskie[i].last_name.substring(0, 1).toLowerCase() == letter) {
				console.log(belorusskie[i].last_name.substring(0, 1).toLowerCase(),letter)
				var name = 
				$(".bel .results").append("<a href='/bel/searchperson?q="+encodeURIComponent(belorusskie[i].last_name+", "+belorusskie[i].full_name)+"'>"+belorusskie[i].last_name+", "+belorusskie[i].full_name+"</a><br>");
			}
		}
}