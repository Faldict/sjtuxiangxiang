function addItem() {
    var number = document.getElementsByClassName('number').length + 1;
    var myForm = document.forms[1];
    var title = myForm["request-title"].value;
    var price = myForm["request-price"].value;
    var perday = myForm["request-days"].value;
    var pertime = myForm["request-times"].value;
    var unit;
    if (perday) {
        unit = price + '元／天';
    } else if (pertime) {
        unit = price + '元/次';
    } else {
        unit = '待商定';
    }
    var xieshang = myForm["request-xieshang"].value;
    var content = myForm["request-content"].value;
    var date1 = document.getElementById('dtp_input1').value;
    var date2 = document.getElementById('dtp_input2').value;
    var item = document.createElement("tr");
    item.innerHTML = '<td><span class="label label-default">未完成</span></td><td class="number">' + number +
                '</td><td>' + title +
                '</td><td>' + date1 + '至' + date2 +
                '</td><td>admin</td><td>' + unit + '</td>';
    document.getElementById('table-body').appendChild(item);
    myForm.reset();
    $('#uploadrequest').modal('hide');
}
