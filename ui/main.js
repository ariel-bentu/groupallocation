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

function loadPupilPrefs() {

    var selTask = getSelectedOption("tasks")
    if (selTask == undefined)
        return
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        $.get("/api/pupil/prefs?task=" + selTask.val() + "&pupilId=" + sel.val(), function (json) {
            emptyList("pupilsPrefs")

            $.each(JSON.parse(json), function (i, value) {
                var name = $("#pupilsInPref1 option[value='" + value.refId + "']").text()
                $('#pupilsPrefs').append($('<option>').text(name).attr('value', value.refId));
            });
        })
            .fail(function () {
                showMessage('Fail in loading pupil prefs');
            });
    }
}



function savePupilPrefs() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        var pupilId = sel.val()
        var pupils = []
        $("#pupilsPrefs option").each(function (i) {
            var opt = $("#pupilsPrefs option").eq(i);
            pupils.push({ "id": pupilId, "refId": opt.val() })
        });
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
            setCurrentByName
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
                "maxAllowed": maxAllowed
            }),
            success: function (result) {
                showMessage("קבוצה עודכנה בהצלחה")
                loadGroups(name);
                $("#saveGroup").prop("disable", true)
            }
        });
    }
}


function
    addPupil(srcList, targetList, btn) {
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

function loadAvailableResults() {
    $.get("/api/available-results?task=" + taskID, function (json) {

        emptyList("results")
        $.each(JSON.parse(json), function (i, value) {
            $('#results').append($('<option>').text(value.runDate + value.title).attr('value', value.id));
        });

        selectFirst("results")
    })
        .done(function () {

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