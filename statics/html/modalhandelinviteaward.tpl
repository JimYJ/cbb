{{define "modalhandelinviteaward"}}
<!--begin::Modal-->
<div class="modal fade" id="AddMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAddRoleTitle" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAddRoleTitle">
                    邀请奖励物品
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/editinvite" enctype="multipart/form-data">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="path" class="form-control-label">
                            选择邀请奖励物品:
                        </label>
                        <select name="itemid" class="form-control m-input--fixed" id="itemid">
                            {{range $v3:=.itemlist}}
                            <option value="{{$v3.id}}">
                                {{$v3.name}}
                            </option>
                            {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="num" class="form-control-label">
                            奖励数量:
                        </label>
                        <input type="text" class="form-control" placeholder="奖励数量" id="num" name="num">
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