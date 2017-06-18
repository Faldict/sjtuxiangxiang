// ----------------------------------------------------
//
// This is TradeRecord of myspace interface javascript.
// Author : Shine
// Date: 2017-06-18
//
// ----------------------------------------------------

function confirmTrade(i) {
	var uri = "/items/shareResponse";
	var uid = document.getElementById('contact' + String(i)).innerHTML;
	var obj_id = document.getElementById('item' + String(i)).innerHTML;
	$.post(uri, {
		"uid_request" : uid,
		"obj_name" : obj_id,
		"agree" : 1
	}, function(status) {
		if (status == '500000') {
			alert("确认成功！");
		} else {
			alert("确认失败！");
		}
	});
}

function getTradeRecord() {
	var url="/user/tradeRecord";
	$.get(url, function(result) {
		var list=document.getElementById("table-body");
		for(var i=0;i<result.length;i++) {
			var item=result[i];
			console.log(item);
			var status1=(item.Cnt==0)?'<td><span class="label label-warning">未完成'
						:'<td><span class="label label-default">已完成';
			var status2=(item.Typ==0)?'<td><span class="label label-primary">租入'
						:'<td><span class="label label-success">租出';
			var status_identify;
			if (item.Cnt == 0 && item.Typ==1) status_identify = '<td><button type="button" class="btn btn-primary btn-sm" onclick="confirmTrade(' +  String(i) + ')">完成交易';
			else if (item.Cnt == 0 && item.Cnt == 0) status_identify = '<td><button type="button" class="btn btn-primary btn-sm" disabled="disabled"">等待确认';
			else status_identify = '<td><button type="button" class="btn btn-sm" disabled="disabled">已经完成';

			var tpl=document.createElement('tr');

				tpl.innerHTML=status1+
				'</span></td>'+status2+
				'</span></td><td><a href="object.html?id='+
				item.Obj_name+'" id="item' + String(i) + '">' + item.Obj_name +
				'</a></td><td><a href="userinformation.html?id=' + item.Uid_other + '" id="contact' + String(i) + '" class="btn btn-link">'+
				item.Uid_other+
				'</button></td><td>'+item.Upload_time+
				'</td>'+status_identify+'</button></td>';
			list.appendChild(tpl);
		}
	});
}

getTradeRecord();
