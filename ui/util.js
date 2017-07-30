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

function emptyList(list) {
    $('#' + list).find('option').remove().end()
}

function selectFirst(list) {
    $('#' + list).val($('#' + list + " option:first").val()).change();
}


function showMessage(msg) {
    $("#message").text( msg);
    setTimeout(clearMessage, 2000);
}

function clearMessage() {
    $("#message").text( "");
}

function searchInList(list, searchStr){ 
    $("#"+ list + " option").each(function(i) {
        var opt = $("#"+ list + " option").eq(i);
        
        if (opt.text().startsWith(searchStr) || opt.text().indexOf(" "+searchStr) >= 0 ) {
            opt.show();
        } else {
            opt.hide();
        }
    });
};
