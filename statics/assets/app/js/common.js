$(document).ready(function () {
    getShowMenu = function (mastermenuid) {
        store.set('nowMenu', mastermenuid)
    }
    initMenu = function () {
        nowMenuID = store.get('nowMenu')
        if (nowMenuID != null) {
            $('li[mid="' + nowMenuID + '"]').each(function () {
                $(this).addClass("m-menu__item--open")
                $(this).find("div").attr("style", "display: block;")
            })
        }
    }
    initMenu()
})