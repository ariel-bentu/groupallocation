

function showResultInfo() {
    var sel = getSelectedOption("results")
    if (sel != undefined) {
        $("#resultName").val(sel.text())
    }
}

function saveResultInfo() {
    var sel = getSelectedOption("results")
    if (sel != undefined) {
        let newName = $("#resultName").val()
        $.ajax({
            url: "/api/result?task=" + taskID + "&id=" + sel.val() + "&resultName=" + decodeURIComponent(newName),
            type: 'POST',
            success: function (result) {
                showMessage("שם עודכן בהצלחה")
                loadAvailableResults(newName)
            }
        });
    }
}

function duplicateResult() {
    var sel = getSelectedOption("results")
    if (sel != undefined) {
        let newName = "copy of " + sel.text()
        let id = sel.val();
        $.ajax({
            url: "/api/result/duplicate?task=" + taskID + "&id=" + sel.val() + "&resultName=" + decodeURIComponent(newName),
            type: 'POST',
            success: function (result) {
                showMessage("שם עודכן בהצלחה")
                loadAvailableResults(newName)
            }
        });
    }
}

function deleteTask() {
    var sel = getSelectedOption("tasks")
    if (sel != undefined) {
        if (window.confirm("האם את בטוחה שאת מעוניינת למחוק את השיבוץ: " + sel.text())) {
            $.ajax({
                url: "/api/tasks?task=" + sel.val(),
                type: 'DELETE',
                success: function (result) {
                    showMessage("קבוצה עודכנה בהצלחה")
                    $("#saveGroup").prop("disable", true)
                }
            });

        }
    }
}

function setDirty(btn) {
    $("#" + btn).prop('disabled', false);
}

function clearDirty(btn) {
    $("#" + btn).prop('disabled', true);
}

function showPupilInfo() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        $("#pupilName").val(sel.data('name'));
        $("#pupilID").val(sel.data('id'));
        $("#pupilGender").val(sel.data('isMale') ? 1 : 2);
        $("#pupilRemarks").val(sel.data('remarks'));
        $("#pupilActive").prop('checked', sel.data('active'));
        clearDirty("btnSavePupil")
    }
}

function loadPupilSubgroups() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        $.get("/api/subgroups/pupil?task=" + taskID + "&pupilId=" + sel.val(), function (json) {
            emptyList("pupilsSubgroups")

            $.each(JSON.parse(json), function (i, value) {
                $('#pupilsSubgroups').append($('<option>').text(value.name));
            });
        })
            .fail(function () {
                showMessage('Fail in loading pupil subgroups');
            });
    }
}

function loadPupilPrefs() {

    var selTask = getSelectedOption("tasks")
    if (selTask == undefined)
        return
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        $.get("/api/pupil/prefs?task=" + selTask.val() + "&pupilId=" + sel.val(), function (json) {
            emptyPrefList()

            $.each(JSON.parse(json), function (i, value) {
                addPrefToList(value, i + 1);
            });
        })
            .fail(function () {
                showMessage('Fail in loading pupil prefs');
            });
    }
}


function emptyPrefList() {
    $('#prefName1').val("");
    $('#prefName2').val("");
    $('#prefName3').val("");
    $('#prefActive1').prop("checked", true);
    $('#prefActive2').prop("checked", true);
    $('#prefActive3').prop("checked", true);
}

function addPrefToList(value, index) {
    var name = $("#pupilsInPref1 option[value='" + value.refId + "']").data('name')

    $('#prefName' + index).val(name);
    $('#prefName' + index).attr('value', value.refId);
    $('#prefActive' + index).prop('checked', value.active);
}

function swapPref(src, dest) {
    let val = $('#prefName' + src).val();
    if (val == "" || $('#prefName' + dest).val() == "") {
        return
    }
    let active = $('#prefActive' + src).prop('checked');

    $('#prefName' + src).val($('#prefName' + dest).val());
    $('#prefActive' + src).prop('checked', $('#prefActive' + dest).prop('checked'));

    $('#prefName' + dest).val(val);
    $('#prefActive' + dest).prop('checked', active);

    setDirty("btnSavePrefs");
}
function addNewPref() {
    if (!getSelectedOption("pupilsInPref1")) {
        return
    }
    var sel = getSelectedOption("pupilsInPref2");
    if (sel != undefined) {
        let id = sel.val();
        var name = $("#pupilsInPref1 option[value='" + id + "']").text()

        for (var j = 1; j < 4; j++) {
            if ($('#prefName' + j).val() == name) {
                showMessage("תלמיד זה נבחר כבר כהעדפה")
                return
            }
        }
        //find empty slot:
        let index = -1
        if ($('#prefName1').val() == "") {
            index = 1
        } else if ($('#prefName2').val() == "") {
            index = 2
        } else if ($('#prefName3').val() == "") {
            index = 3
        }
        if (index < 0) {
            showMessage("כל ההעדפות נבחרו")
            return;
        }

        $('#prefName' + index).val(name);
        $('#prefName' + index).attr('value', id);

        setDirty("btnSavePrefs")
    }

}

function removePrefFromList(index) {
    for (var i = index; i < 3; i++) {
        let j = i + 1
        $('#prefName' + i).val($('#prefName' + j).val());
        $('#prefName' + i).attr('value', $('#prefName' + j).attr('value'));
        $('#prefName' + j).val("");
        $('#prefName' + j).attr('value', "");
    }
    if (index == 3) {
        $('#prefName3').val("");
        $('#prefName3').attr('value', "");
    }
    setDirty("btnSavePrefs");
}

function savePupilPrefs() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        var pupilId = sel.val()
        var pupils = []
        for (var i = 1; i < 4; i++) {
            if ($("#prefName" + i).val() != "") {
                pupils.push({
                    "id": pupilId,
                    "refId": $("#prefName" + i).attr('value'),
                    "active": $("#prefActive" + i).prop('checked')
                })
            }
        }

        var jsonStr = JSON.stringify(pupils);

        $.post("/api/pupil/prefs?task=" + taskID + "&pupilId=" + pupilId, jsonStr, function (json) {
            showMessage('Successfully saved');
            $("#btnSavePrefs").prop('disabled', true);

        })
            .done(function () {

            })
            .fail(function () {
                showMessage("Fail in saving pupil's prefs");
            })
            .always(function () {
            });

    }

}

function savePupil() {
    //validate
    if ($("#pupilGender").val() != 1 && $("#pupilGender").val() != 2) {
        showMessage('חובה לבחור מין')
        return;
    }
    var ID = $("#pupilID").val();
    var name = $("#pupilName").val();
    var isMale = $("#pupilGender").val() == 1;
    var remarks = $("#pupilRemarks").val();
    var active = $("#pupilActive").prop('checked');

    var url = "/api/pupils?task=" + taskID + "&ID=" + ID

    $.ajax({
        url: url,
        type: ID == "" ? 'PUT' : 'POST',
        data: JSON.stringify({
            "name": name,
            "remarks": remarks,
            "isMale": isMale,
            "active": active
        }),
        success: function (result) {
            showMessage("תלמיד נשמר בהצלחה")
            loadPupils(name);

            clearDirty("btnSavePupil")
        }
    });

}
function deletePupil() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        if (window.confirm("האם את בטוחה שאת מעוניינת למחוק את התלמיד: " + sel.text())) {
            var ID = $("#pupilID").val()

            var url = "/api/pupils?task=" + taskID + "&id=" + ID
            $.ajax({
                url: url,
                type: 'DELETE',
                success: function (result) {
                    showMessage("תלמיד נמחק בהצלחה")
                    loadPupils();
                }
            });
        }
    }
}

function deleteSubgroup() {
    var sel = getSelectedOption("subgroups");
    if (sel != undefined) {
        if (window.confirm("האם את בטוחה שאת מעוניינת למחוק את הקבוצה: " + sel.text())) {
            var ID = $("#groupID").val()

            var url = "/api/subgroup?task=" + taskID + "&groupId=" + ID
            $.ajax({
                url: url,
                type: 'DELETE',
                success: function (result) {
                    showMessage("קבוצה נמחקה בהצלחה")
                    loadGroups();
                }
            });
        }
    }
}
function saveSubgroup() {
    var sel = getSelectedOption("subgroups");
    if (sel != undefined) {
        var ID = $("#groupID").val()
        var name = $("#groupName").val()
        var isUnite = $("#isUnite").prop('checked');
        var isInactive = $("#isInactive").prop('checked');
        var isSpreadEvenly = $("#isSpreadEvenly").prop('checked');
        var isGenderSensitive = $("#isGenderSensitive").prop('checked');
        var minAllowed = $("#minAllowed").val();
        var maxAllowed = $("#maxAllowed").val();
        var isGarden = $("#isGarden").prop('checked')

        var url = "/api/subgroup?task=" + taskID + "&groupId=" + ID
        $.ajax({
            url: url,
            type: ID == "" ? 'PUT' : 'POST',
            data: JSON.stringify({
                "name": name,
                "isUnite": isUnite,
                "isSpreadEvenly": isSpreadEvenly,
                "isGenderSensitive": isGenderSensitive,
                "isInactive": isInactive,
                "minAllowed": minAllowed,
                "maxAllowed": maxAllowed,
                "isGarden": isGarden
            }),
            success: function (result) {
                showMessage("קבוצה עודכנה בהצלחה")
                loadGroups(name);
                $("#saveGroup").prop("disable", true)
            }
        });
    }
}


function addPupil(srcList, targetList, btn) {
    var sel = getSelectedOption(srcList);
    if (sel != undefined) {
        addOptionToList(targetList, sel)
        setDirty(btn)
    }
}

function removePupil(srcList, btn) {
    var sel = getSelectedOption(srcList);
    if (sel != undefined) {
        removeOptionFromList(srcList, sel)
        setDirty(btn)
    }
}

function loadPupils(currentName) {
    let currentID = undefined;
    $.get("/api/pupils?task=" + taskID, function (json) {
        emptyList("pupils")
        emptyList("pupilsInPref1")
        emptyList("pupilsInPref2")

        $.each(JSON.parse(json), function (i, value) {
            var name = decrypt($("#pwd").val(), value.name)
            var opt = $('<option>').text(name + " - " + (value.isMale ? "בן" : "בת"))
            opt.attr('value', value.id);
            $(opt).data('name', name);
            $(opt).data('active', value.active);
            $(opt).data('isMale', value.isMale);
            $(opt).data('remarks', value.remarks);
            $(opt).data('id', value.id);

            if (currentName && currentName == name) {
                currentID = value.id
            }

            $('#pupilsInPref1').append(opt);
            $('#pupilsInPref2').append($('<option>').text(name + " - " + (value.isMale ? "בן" : "בת")).attr('value', value.id));
            $('#pupils').append($('<option>').text(name + " - " + (value.isMale ? "בן" : "בת")).attr('value', value.id));
        });
    })
        .done(function () {
            if (currentID) {
                $("#pupilsInPref1").val(currentID).change();
                $("#pupilsInPref1").val(currentID).click();
            }
        })
        .fail(function () {
            showMessage('Fail in loading pupils');
        })
        .always(function () {
        });
}

function loadGroups(currentName) {
    let currentID = undefined;

    $.get("/api/subgroups?task=" + taskID, function (json) {
        emptyList("subgroups");
        $.each(JSON.parse(json), function (i, value) {
            var opt = $('<option>').text(value.name + " - \t\t" + (value.isUnite ? "איחוד" : "פירוד")).attr('value', value.id);
            $(opt).data('id', value.id)
            $(opt).data('name', value.name)
            $(opt).data('isUnite', value.isUnite)
            $(opt).data('isInactive', value.isInactive);
            $(opt).data('isSpreadEvenly', value.isSpreadEvenly);
            $(opt).data('isGenderSensitive', value.isGenderSensitive);
            $(opt).data('minAllowed', value.minAllowed);
            $(opt).data('maxAllowed', value.maxAllowed);
            $(opt).data('isGarden', value.isGarden);

            if (currentName && currentName == value.name) {
                currentID = value.id;
            }
            $('#subgroups').append(opt)
        });
        selectFirst("subgroups")
        //loadSubgroupPupils();
    })
        .done(function () {
            if (currentID) {
                $('#subgroups').val(currentID).change();
                $('#subgroups').val(currentID).click();
            }
        })
        .fail(function () {
            showMessage('Fail in loading pupils');
        })
        .always(function () {
        });
}

function showGroupInfo() {
    var sel = getSelectedOption("subgroups")
    if (sel != undefined) {
        $("#groupID").val(sel.data('id'));
        $("#groupName").val(sel.data('name'));
        $("#isUnite").prop('checked', sel.data('isUnite'));
        $("#isGenderSensitive").prop('checked', sel.data('isGenderSensitive'));
        $("#isSpreadEvenly").prop('checked', sel.data('isSpreadEvenly'));
        $("#isInactive").prop('checked', sel.data('isInactive'));
        $("#isGarden").prop('checked', sel.data('isGarden'));
        $("#maxAllowed").val(sel.data('maxAllowed'));
        $("#minAllowed").val(sel.data('minAllowed'));
    }
}

function loadSubgroupPupils() {
    var sel = getSelectedOption("subgroups")
    emptyList("groupsPupils");
    if (sel != undefined) {
        $.get("/api/subgroup/pupils?task=" + taskID + "&groupId=" + sel.val(), function (json) {

            $.each(JSON.parse(json), function (i, value) {
                var name = $("#pupils option[value='" + value.refId + "']").text()
                $('#groupsPupils').append($('<option>').text(name).attr('value', value.refId));
            });


        })
            .done(function () {

            })
            .fail(function () {
                showMessage('Fail in loading pupils');
            })
            .always(function () {
            });
    }
}

function loadTasks() {
    $.get("/api/tasks", function (json) {

        emptyList("tasks")
        $.each(JSON.parse(json), function (i, value) {
            $('#tasks').append($('<option>').text(value.name).attr('value', value.id));
        });

        selectFirst("tasks")
    })
        .done(function () {

        })
        .fail(function () {
            showMessage('Fail in loading tasks');
        })
        .always(function () {
        });
}

function loadAvailableResults(currentName) {
    let currentId = undefined
    $.get("/api/available-results?task=" + taskID, function (json) {

        emptyList("results")
        $.each(JSON.parse(json), function (i, value) {
            if (currentName && currentName == value.title) {
                currentId = value.id
            }
            let title = value.title != ""? value.title : value.runDate;
            $('#results').append($('<option>').text(title).attr('value', value.id));
        });

        
    })
        .done(function () {
            if (currentId) {
                $("#results").val(currentId).change();
                $("#results").val(currentId).click();
            } else {
                selectFirst("results")
            }

        })
        .fail(function () {
            showMessage('Fail in loading tasks');
        })
        .always(function () {
        });
}

function saveGroupsPupils() {

    var sel = getSelectedOption("subgroups")
    if (sel != undefined) {
        var groupId = sel.val()
        var pupils = []
        $("#groupsPupils option").each(function (i) {
            var opt = $("#groupsPupils option").eq(i);
            pupils.push({ "id": groupId, "refId": opt.val() })
        });
        var jsonStr = JSON.stringify(pupils);

        $.post("/api/subgroup/pupils?task=" + taskID + "&groupId=" + groupId, jsonStr, function (json) {
            showMessage('Successfully saved');

        })
            .done(function () {

            })
            .fail(function () {
                showMessage('Fail in loading pupils');
            })
            .always(function () {
            });

    }
}


function showMessage(msg) {
    $('#message').text(msg);
    setTimeout(() => $('#message').text(""), 3000);
}