{{define "modalHandelTreeLevel"}}
<!--begin::Modal-->
<div class="modal fade" id="AddMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAddRoleTitle" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAddRoleTitle">
                    编辑桑树等级
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/edittreelevel">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="redeemitem" class="form-control-label">
                            每小时产桑叶数量:
                        </label>
                        <input type="text" class="form-control" placeholder="每小时产桑叶数量" id="growthhours" name="growthhours">
                    </div>
                    <div class="form-group">
                        <label for="redeemitem" class="form-control-label">
                            最大可积累时间(小时):
                        </label>
                        <input type="text" class="form-control" placeholder="最大可积累时间(小时)" id="maxhours" name="maxhours">
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