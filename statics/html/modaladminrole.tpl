{{define "modalAdminRole"}}
<!--begin::Modal-->
<div class="modal fade" id="AdminRole" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    设置用户岗位
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="roleAction" action="/adminbindrole">
                <div class="modal-body">
                    <div class="form-group m-form__group">
                        <label for="roleSelect">
                            所属岗位:(可多选)
                        </label>
                        <select multiple="multiple" class="form-control m-input" id="roleSelect" name="roles">
                            {{range $r:= .rolelist}}
                            <option value="{{$r.id}}">
                                    {{$r.name}}
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