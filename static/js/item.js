// ---------------------------------------
//
// This is the item interface javascript.
// Author : Faldict
// Date: 2017-06-01
//
// ---------------------------------------

// String.format
// seemed to be useless
String.prototype.format = function() {
  var args = arguments;
  return this.replace(/\{(\d+)\}/g,function(s,i){
    return args[i];
  });
}

// user add an item
function addItem() {
    // define segments
    var number = document.getElementsByClassName('number').length + 1;
    var myForm = document.forms[1];
    var title = myForm["request-title"].value;
    var price = myForm["request-price"].value;
    var perday = myForm["request-days"].value;
    var pertime = myForm["request-times"].value;
    var content = myForm["request-content"].value;
    var date1 = document.getElementById('dtp_input1').value;
    var date2 = document.getElementById('dtp_input2').value;
    // hide modal
    $('#uploadrequest').modal('hide');
    // call api via ajax
    var uri = "/items/add";
    $.post(uri, {
            "obj_name" : title,
            "obj_price" : price,
            "info" : content,
            // "use_time" : date1,
            "start_time" : date1,
            "end_time" : date2
        }, function (status) {
            if (status == "400000") {
                alert("Upload Success!");
            } else {
                alert("Upload Failed! Try it later!");
            }
        }
    );
}

// Get Item List
function getRentItem() {
    var url = "items/listItem?type=1";
    $.getJSON(url, function(result) {
        var list = document.getElementById("table-body");
        for (var i=0; i<result.length; i++) {
            var item = result[i];
            var tpl = document.createElement('div');
            tpl.innerHTML = '<div class="thumbnail"><img src="http://www1.pcbaby.com.cn/yongping624.jpg" alt="..."><div class="caption"><a href="object.html?id='
                + item.obj_name + '"><h4>' + item.obj_name +
                '</h4></a><p>' + item.obj_info +
                '</p><p class="text-muted text-right"><span class="label label-success">' +
                item.obj_state + '</span>' + item.obj_price + '元/天</p></div></div>';
            tpl.class = "col-md-4";
            list.appendChild(tpl);
        }
    });
}

function getWantItem() {
    var url = "items/listItem?type=0";
    $.get(url, function(result) {
        var list = document.getElementById("table-body");
        for (var i=0; i<result.length; i++) {
            var item = result[i];
            console.log(item);
            var status = item.Obj_state == "0" ? '<td><span class="label label-default">未完成' : '<td><span class="label label-success">已完成';
            var tpl = document.createElement('tr');
                tpl.innerHTML = status +
                '</span></td><td><a href="object.html?id=' +
                item.Obj_name + '">' + item.Obj_name +
                '</a></td><td>' + item.Upload_time + '至' + item.Use_time + '</td><td>' + item.Uid + '</td><td>' + item.Obj_price + '元/天 </td>';
            // tpl.class = "success";
            list.appendChild(tpl);
        }
    })
}
