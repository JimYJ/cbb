
//== Class definition

var Menu = function () {
  //== Private functions

  // demo initializer
  var option = function () {

    var datatable = $('.m-datatable').mDatatable({
      data: {
        saveState: { cookie: false },
      },
      search: {
        input: $('#generalSearch'),
      },
      pageSize: 15,
      serverPaging: true,
      columns: [
        {
          field: 'id',
          title: "id",
          width: 10,
        },
        {
          field: 'name',
          width: 50,
        },
        {
          field: '所属父级',
          width: 50,
        },
        {
          field: '路径',
          width: 10,
        },
        {
          field: '创建时间',
          type: 'date',
          format: 'YYYY-MM-DD',
          width: 80,
        },
        {
          field: '最后编辑',
          type: 'date',
          format: 'YYYY-MM-DD',
          width: 80,
        },
        {
          field: '操作',
          width: 80,
        },
      ],
    });
  };

  return {
    //== Public functions
    init: function () {
      // init dmeo
      option();
    },
  };
}();


//== Class definition
var Select2 = function () {
  //== Private functions
  var demos = function () {

    // multi select
    $('#menuSelect, #m_select2_3_validate').select2({
      placeholder: "选择可访问菜单",
      width: "100%",
      allowClear: true,
      tags: true,
      multiple: true,
      // maximumSelectionLength : 3,
      selectOnBlur:true,
      // formatSelectionTooBig:"你只能选中三个"
    });
  }


  //== Public functions
  return {
    init: function () {
      demos();
    }
  };
}();

jQuery(document).ready(function () {
  Menu.init();
  Select2.init();
});