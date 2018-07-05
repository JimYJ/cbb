{{define "modalHandelVendor"}}
<!--begin::Modal-->
<div class="modal fade" id="AddEdit" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    新增店铺
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/addvendor">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="names" class="form-control-label">
                            店铺名称:
                        </label>
                        <input type="text" class="form-control" placeholder="店铺名称" id="name" name="names">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            负责人:
                        </label>
                        <input type="text" class="form-control" placeholder="负责人" id="leader" name="leader">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            负责人电话:
                        </label>
                        <input type="text" class="form-control" placeholder="负责人电话" id="leaderphone" name="leaderphone">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            省市县选择:
                        </label>
                        <select id="prov" onchange="showCity(this)" class="form-control m-input--fixed" name="province">
                        </select>
                        <!--城市选择-->
                        <select id="city" onchange="showCountry(this)" class="form-control m-input--fixed" name="city">
                        </select>
                        <!--县区选择-->
                        <select id="country" onchange="selecCountry(this)" class="form-control m-input--fixed" name="county">
                        </select>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">
                        取消
                    </button>
                    <input value="确定" type="submit" id="btns" class="btn btn-primary">
                    </input >
                </div>
            </form>
        </div>
    </div>
</div>
<!--end::Modal-->
{{end}}