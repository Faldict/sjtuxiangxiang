// ---------------------------------------
//
// This is the item interface javascript.
// Author : Faldict
// Date: 2017-06-01
//
// ---------------------------------------

// String.format
String.prototype.format = function() {
  var args = arguments;
  return this.replace(/\{(\d+)\}/g,function(s,i){
    return args[i];
  });
}

// user add an item
function addItem() {

}

// Get Item List
function getItem() {
    var url = "items/listItem";
    $.get(url, function(result) {
        var list = document.getElementById("itemlist");
        for (var i=0; i<result.length; i++) {
            var item = result[i];
            var tpl = document.createElement;
            tpl.innerHTML = '<div class="thumbnail"><img src="http://www1.pcbaby.com.cn/yongping624.jpg" alt="..."><div class="caption"><a href="#"><h4>' + item.obj_name +
                '</h4></a><p>' + item.obj_info +
                '</p><p class="text-muted text-right"><span class="label label-success">' +
                item.obj_state + '</span>' + item.obj_price + '元/天</p></div></div>';
            tpl.class = "col-md-4";
            list.appendChild(tpl);
        }
    })
}
