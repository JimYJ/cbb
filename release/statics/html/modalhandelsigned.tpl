{{define "modalHandelSigned"}}
<!--begin::Modal-->
<div class="modal fade" id="AddMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAddRoleTitle" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAddRoleTitle">
                    编辑签到规则
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/editsigned" enctype="multipart/form-data">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            签到简介
                        </label>
                        <textarea id="intro" name="intro"></textarea>
                    </div>
                    <div class="form-group">
                        <label for="path" class="form-control-label">
                            选择每日签到奖励:
                        </label>
                        <select name="dayitemid" class="form-control m-input--fixed" id="dayitemid">
                            {{range $v3:=.itemlist}}
                            <option value="{{$v3.id}}">
                                {{$v3.name}}
                            </option>
                            {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                            <label for="path" class="form-control-label">
                                选择整周签到奖励:
                            </label>
                            <select name="weekitemid" class="form-control m-input--fixed" id="weekitemid">
                                {{range $v4:=.itemlist}}
                                <option value="{{$v4.id}}">
                                    {{$v4.name}}
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
                    </input>
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}