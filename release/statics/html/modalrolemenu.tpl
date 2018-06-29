{{define "modalRoleMenu"}}
<!--begin::Modal-->
<div class="modal fade" id="RoleMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    设置岗位权限
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="roleAction" action="/rolebindmenu">
                <div class="modal-body">
                    <div class="form-group m-form__group row">
                        <label class="col-form-label col-lg-12 col-sm-12">
                            选择可访问页面
                        </label>
                        <div class="col-lg-12 col-md-12 col-sm-12">
                            <select class="form-control m-select2" id="menuSelect" name="rolemenu" multiple="multiple">
                                {{range $m:= .menulist}}
                                <optgroup label="{{$m.name}}">
                                    {{range $s:= $m.list}}
                                    <option value="{{$s.id}}">
                                        {{$s.name}}
                                    </option>
                                    {{end}}
                                </optgroup>
                                {{end}}
                            </select>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">
                        取消
                    </button>
                    <input value="确定" type="submit" class="btn btn-primary">
                    </input>
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}