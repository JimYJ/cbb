{{define "modalHandelQuestion"}}
<!--begin::Modal-->
<div class="modal fade" id="AddEdit" tabindex="-1" role="dialog" aria-labelledby="ModalAdd" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="ModalAdd">
                    新增记录
                </h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">
                        &times;
                    </span>
                </button>
            </div>
            <form method="POST" id="handelAction" action="/addquestion">
                <div class="modal-body">
                    <div class="form-group">
                        <label for="names" class="form-control-label">
                            问题内容:
                        </label>
                        <input type="text" class="form-control" placeholder="问题内容" id="content" name="content">
                    </div>
                    <div class="form-group">
                        <label for="name" class="form-control-label">
                            兑换物品:
                        </label>
                        <select name="itemid" class="form-control m-input--fixed" id="itemid">
                            {{range $c:=.itemlist}}
                            <option value="{{$c.id}}">
                                {{$c.name}}
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