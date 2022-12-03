const navId = "profile-container";
const menuId = "profile-menu"

$("body").click(function (e) {
    if (e.target.id === menuId) return;
    if ($(e.target).closest(`#${menuId}`).length > 0) return;
    $(`#${navId}`).removeClass("open");
    $(`#${menuId}`).removeClass("open");
})

$(`#${navId}`).on("click", function (e) {
    $(`#${navId}`).toggleClass("open");
    $(`#${menuId}`).toggleClass("open");
    e.stopPropagation();
})