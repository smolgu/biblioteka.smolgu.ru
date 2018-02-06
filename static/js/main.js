// Simple JavaScript Templating
// John Resig - http://ejohn.org/ - MIT Licensed
(function(){
  var cache = {};
 
  this.tmpl = function tmpl(str, data){
    // Figure out if we're getting a template, or if we need to
    // load the template - and be sure to cache the result.
    var fn = !/\W/.test(str) ?
      cache[str] = cache[str] ||
        tmpl(document.getElementById(str).innerHTML) :
     
      // Generate a reusable function that will serve as a template
      // generator (and which will be cached).
      new Function("obj",
        "var p=[],print=function(){p.push.apply(p,arguments);};" +
       
        // Introduce the data as local variables using with(){}
        "with(obj){p.push('" +
       
        // Convert the template into pure JavaScript
        str
          .replace(/[\r\t\n]/g, " ")
          .split("<%").join("\t")
          .replace(/((^|%>)[^\t]*)'/g, "$1\r")
          .replace(/\t=(.*?)%>/g, "',$1,'")
          .split("\t").join("');")
          .split("%>").join("p.push('")
          .split("\r").join("\\'")
      + "');}return p.join('');");
   
    // Provide some basic currying to the user
    return data ? fn( data ) : fn;
  };
})();


function init() {
 $(function() {
        $(".nav-left").accordion({
            icons: null,
            heightStyle: "content"
        });
    });
    $(function() {
        $(".link-accordion").accordion({
            icons: null,
            heightStyle: "content"
        });
    });
    $(".fancybox").fancybox({
        prevEffect: 'elastic',
        nextEffect: 'elastic',
        padding: 0,
        closeBtn: true,
        helpers: {
            title: {
                type: 'inside'
            },
        }
    });
    $(".prezi").fancybox({
        prevEffect: 'elastic',
        nextEffect: 'elastic',
        padding: 0,
        closeBtn: true,
        helpers: {
            title: {
                type: 'inside'
            },
            buttons  : {}
        }
    });
    $("img.lazy").lazyload();

if (location.host === "biblioteka.smolgu.ru") {
new Image().src = "//counter.yadro.ru/hit?r" +
    escape(document.referrer) + ((typeof(screen) == "undefined") ? "" :
        ";s" + screen.width + "*" + screen.height + "*" + (screen.colorDepth ?
        screen.colorDepth : screen.pixelDepth)) + ";u" + escape(document.URL) +
    ";" + Math.random();
}



(function(a){a.fn.removeWhitespace=function(){this.contents().filter(function(){return this.nodeType==3&&!/\S/.test(this.nodeValue)}).remove();return this}})(jQuery)


/*$(window).load(function(){
    $('.collage').collagePlus({
     'targetHeight'    : 100,
 });
})*/
 




   function scrollTop() {
    return (window.pageYOffset !== undefined) ? window.pageYOffset : (document.documentElement || document.body.parentNode || document.body).scrollTop;
  }

  window.addEventListener("scroll",function(){
        var yPos = -(scrollTop() / 2); 
        var coords = 'center '+ yPos + 'px';
        
        document.body.style.backgroundPosition =coords
  });

}

$(document).ready(function() {
   init()

var summer_bgs = ["leto.jpg", "leto-2.jpg", "leto-1.jpg"],
autumn_bgs = ["window/autumn/0.jpg","window/autumn/1.jpg","window/autumn/2.jpg","window/autumn/3.jpg","window/autumn/4.jpg"],
winter_bgs = ["window/winter/0.jpg","window/winter/1.jpg","window/winter/2.jpg","window/winter/3.jpg","window/winter/4.jpg","window/winter/5.jpg"],
current = 0,
bgs = summer_bgs, mn

var d = new Date();
mn = d.getMonth()

if (mn >= 8) {
    bgs = autumn_bgs
}
if (mn == 0 || mn == 1 || mn == 11) {
  bgs = winter_bgs
  $(".header-gradient").addClass("snow");
  $("a").addClass("snow-hover");
}

var el = $(".header-window")

function changeBg() {
    current++
    if (current > (bgs.length-1)) {
        current = 0
    }
    el.css("background-image", "url('/static/img/" +
        bgs[current] + "')");
}
changeBg();
setInterval(changeBg, 14000)

el.css({opacity:1});
});

$(window).on("statechangecomplete", function() {
    init()
})

loadCSS("/"+smolgu_lib.build_hash+"/fancybox/jquery.fancybox.css");
loadCSS("/"+smolgu_lib.build_hash+"/fancybox/helpers/jquery.fancybox-buttons.css");

/* показываем картинку на переход к статистике ЛИ */

$("#li-counter").html(`<a href="http://www.liveinternet.ru/click" target="_blank"><img src="//counter.yadro.ru/logo?14.5" title="LiveInternet: показано число просмотров за 24 часа, посетителей за 24 часа и за сегодня" alt="" border="0" width="88" height="31"/></a>`);