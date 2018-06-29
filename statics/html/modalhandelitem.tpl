{{define "modalHandelItem"}}
<!--begin::Modal-->
<div class="modal fade" id="AddMenu" tabindex="-1" role="dialog" aria-labelledby="ModalAddRoleTitle" aria-hidden="true">
    <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAddRoleTitle">
                    编辑物品
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/edititem">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            增加经验值:
                        </label>
                        <input type="text" class="form-control" placeholder="增加经验值" id="exp" name="exp">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            每日使用限制:
                        </label>
                        <input type="text" class="form-control" placeholder="每日使用限制" id="limitday" name="limitday">
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