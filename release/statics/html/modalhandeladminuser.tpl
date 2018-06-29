{{define "modalHandelAdminUser"}}
<!--begin::Modal-->
<div class="modal fade" id="AddEdit" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    新增用户
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/addadmin">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            用户:
                        </label>
                        <input type="text" class="form-control" placeholder="用户" id="names" name="username">
                    </div>
                    <div class="form-group">
                        <label for="pass" class="form-control-label">
                            密码:
                        </label>
                        <input type="password" class="form-control" placeholder="密码" id="pass" name="password">
                    </div>
                    <div class="form-group">
                        <label for="parentid" class="form-control-label">
                            选择所属店铺
                        </label>
                        <select name="vendorid" class="form-control m-input--fixed" id="vendorid">
                            <option value="0">
                                [平台用户]
                            </option>
                            {{range $v3:=.vendorlist}}
                            <option value="{{$v3.id}}">
                                {{$v3.name}}
                            </option>
                            {{end}}
                        </select>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">
                        取消
                    </button>
                    <input value="确定" type="submit" class="btn btn-primary">
                    </input >
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}