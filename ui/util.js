function isSelected(list) {
    
    return getSelectedOption(list) != undefined;
}

function addListItem(list, itemVal, itemText) {
    if (!isOptionExists(list, itemVal))
        $("#" + list).append("<option value=\""+ itemVal + "\">"+ itemText +"</option>")
    else 
        showMessage("Already selected")
}

function removeListItem(list, itemVal) {
    $("#"+ list +" option[value=\"" + itemVal + "\"]").remove();
}

function getSelectedOption(list) {
    var sel = $("#" + list + " option:selected" );
    if (sel.val() != undefined) {
        return sel;
    }
    return undefined
}

function addOptionToList(list, opt) {
    addListItem(list, opt.val(), opt.text())
}

function removeOptionFromList(list, opt) {
    removeListItem(list, opt.val());
}


function isOptionExists(list, optVal) {
    return ($("#"+ list +" option[value=\"" + optVal + "\"]").val() != undefined)
}

function showMessage(msg) {
    $("#message").text( msg);
    setTimeout(clearMessage, 2000);
}

function clearMessage() {
    $("#message").text( "");
}

