{{define "modalHandelRole"}}
<!--begin::Modal-->
<div class="modal fade" id="AddMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAddRoleTitle" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAddRoleTitle">
                    新增角色
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/addrole">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            名称:
                        </label>
                        <input type="text" class="form-control" placeholder="名称" id="names" name="name">
                    </div>
                    <!-- <div class="form-group">
                        <label for="parentid" class="form-control-label">
                            选择所属父级:(例：/test)
                        </label>
                        <select name="parentid" class="form-control m-input--fixed" id="parentid">
                            <option value="0">
                                [父级]
                            </option>
                            {{range $v3:=.mainlist}}
                            <option value="{{$v3.id}}">
                                {{$v3.name}}
                            </option>
                            {{end}}
                        </select>
                    </div> -->
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