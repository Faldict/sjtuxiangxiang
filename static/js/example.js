function addItem() {
    var number = document.getElementsByClassName('number').length + 1;
    var item = document.createElement("tr");
    item.innerHTML = '<td><span class="label label-default">未完成</span></td><td class="number">' + number +
    '</td><td>求借单反相机进行社团活动拍摄</td><td>2017.1.1~2017.1.2</td><td>admin</td><td>50元/天</td>';
    document.getElementById('table-body').appendChild(item);
}
