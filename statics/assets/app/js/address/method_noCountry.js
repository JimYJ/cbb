
/*根据id获取对象*/
function $A(str) {
    return document.getElementById(str);
}

var addrShow = $A('addr-show');
var btn = $A("btns");
var prov = $A('prov');
var city = $A('city');
var country = $A('country');


/*用于保存当前所选的省市区*/
var current = {
    prov: '',
    city: '',
    country: ''
};

/*自动加载省份列表*/
function showProv() {
    // btn.disabled = true;
    var len = provice.length;
    for (var i = 0; i < len; i++) {
        // $("#prov").append("<option value='"+i+"'>"+provice[i]['name']+"</option>");
        var provOpt = document.createElement('option');
        provOpt.innerText = provice[i]['name'];
        provOpt.value = i;
        prov.appendChild(provOpt);
    }
};

/*根据所选的省份来显示城市列表*/
function showCity(obj) {
    var val = obj.options[obj.selectedIndex].value;
    if (val != current.prov) {
        current.prov = val;
        // addrShow.value = '';
        // btn.disabled = true;
    }
    //console.log(val);
    if (val != null) {
        city.length = 1;
        var cityLen = provice[val]["city"].length;
        for (var j = 0; j < cityLen; j++) {
            var cityOpt = document.createElement('option');
            cityOpt.innerText = provice[val]["city"][j].name;
            cityOpt.value = j;
            city.appendChild(cityOpt);
        }
    }
}


/*选择县区之后的处理函数*/
function selecCountry(obj) {
    current.country = obj.options[obj.selectedIndex].value;
    if ((current.city != null) && (current.country != null)) {
        btn.disabled = false;
    }
}