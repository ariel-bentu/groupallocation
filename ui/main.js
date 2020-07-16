function deleteTask() {
    var sel = getSelectedOption("tasks")
    if (sel != undefined) {
        if (window.confirm("האם את בטוחה שאת מעוניינת למחוק את השיבוץ: " + sel.text())) {
            $.ajax({
                url: "/api/tasks?task=" + sel.val(),
                type: 'DELETE',
                success: function(result) {
                    showMessage("קבוצה עודכנה בהצלחה")
                    $("#saveGroup").prop("disable", true)
                }
            });
    
        }
    }
}

function setDirty(btn) {
    $("#"+btn).prop('disabled', false);
}

function loadPupilPrefs() {
    
    var selTask = getSelectedOption("tasks")
    if (selTask == undefined) 
        return
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
    $.get("/api/pupil/prefs?task="+selTask.val()+"&pupilId="+sel.val() , function(json) {
        emptyList("pupilsPrefs")

        $.each(JSON.parse(json), function(i, value) {
            var name = $("#pupilsInPref1 option[value='"+ value.refId + "']").text()
            $('#pupilsPrefs').append($('<option>').text(name).attr('value', value.refId));
        });
    })
    .fail(function() {
        showMessage('Fail in loading pupil prefs');
    });        
    }    
}

function savePupilPrefs() {
    var sel = getSelectedOption("pupilsInPref1");
    if (sel != undefined) {
        var pupilId = sel.val()
        var pupils = []
        $("#pupilsPrefs option").each(function(i) {
            var opt = $("#pupilsPrefs option").eq(i);
            pupils.push({"id":pupilId, "refId":opt.val()})
        });
        var jsonStr = JSON.stringify(pupils);
        
        $.post("/api/pupil/prefs?task="+taskID+"&pupilId="+ pupilId, jsonStr, function(json) {
            showMessage('Successfully saved');
            $("#btnSavePrefs").prop('disabled', true);

        })
        .done(function() {
            
        })
        .fail(function() {
            showMessage("Fail in saving pupil's prefs");
        })
        .always(function() {
        });
        
    }
    
}

function saveGroup() {
    var sel = getSelectedOption("subgroups");
    if (sel != undefined) {
        
        var isUnite = $("#isUnite").prop('checked');
        var isInactive = $("#isInactive").prop('checked');
        var isSpreadEvenly = $("#isSpreadEvenly").prop('checked');
        var isGenderSensitive = $("#isGenderSensitive").prop('checked');
        var minAllowed = $("#minAllowed").val();
        var maxAllowed = $("#maxAllowed").val();
        
        var url = "/api/subgroup?task="+taskID+"&groupId="+sel.val() + "&isUnite="+isUnite+"&isInactive="+isInactive+"&isSpreadEvenly="+isSpreadEvenly+"&isGenderSensitive="+isGenderSensitive+
          "&minAllowed="+minAllowed+"&maxAllowed="+maxAllowed
        $.ajax({
            url: url,
            type: 'PUT',
            success: function(result) {
                showMessage("קבוצה עודכנה בהצלחה")
                loadGroups();
                $("#saveGroup").prop("disable", true)
            }
        });
    }
}

function saveNewGroup() {
    
    $.post("/api/subgroup?task="+taskID + "&name="+$("#mdlGroupName").val(), "", function(json) {
        showMessage('Successfully saved');
        loadGroups();
    })
    .done(function() {
        
    })
    .fail(function() {
        showMessage('Fail saving new group');
    })
    .always(function() {
    });
    $("#newSubgroupModal").hide();
    
}

function addPupil(srcList, targetList, btn){
    var sel = getSelectedOption(srcList);
    if (sel != undefined) {
        addOptionToList(targetList, sel)
        setDirty(btn)
    }
}

function removePupil(srcList, btn){
    var sel = getSelectedOption(srcList);
    if (sel != undefined) {
        removeOptionFromList(srcList, sel)
        setDirty(btn)
    }
}

function loadPupils() {
    $.get("/api/pupils?task="+taskID , function(json) {
        emptyList("pupils")
        emptyList("pupilsInPref1")
        emptyList("pupilsInPref2")

        $.each(JSON.parse(json), function(i, value) {
            var name = decrypt($("#pwd").val(), value.name)

            $('#pupils').append($('<option>').text(name + " - " + (value.isMale?"בן":"בת")).attr('value', value.id));
            $('#pupilsInPref1').append($('<option>').text(name + " - " + (value.isMale?"בן":"בת")).attr('value', value.id));
            $('#pupilsInPref2').append($('<option>').text(name + " - " + (value.isMale?"בן":"בת")).attr('value', value.id));
        });
    })
    .done(function() {
        
    })
    .fail(function() {
        showMessage('Fail in loading pupils');
    })
    .always(function() {
    });
}

function loadGroups() {
    $.get("/api/subgroups?task="+taskID , function(json) {
        emptyList("subgroups");
        $.each(JSON.parse(json), function(i, value) {
            var opt = $('<option>').text(value.name + " - \t\t"+(value.isUnite?"איחוד":"פירוד")).attr('value', value.id);
            $(opt).data('isUnite',value.isUnite)
            $(opt).data('isInactive',value.isInactive);
            $(opt).data('isSpreadEvenly',value.isSpreadEvenly);
            $(opt).data('isGenderSensitive',value.isGenderSensitive);
            $(opt).data('minAllowed',value.minAllowed);
            $(opt).data('maxAllowed',value.maxAllowed);

            $('#subgroups').append(opt)
        });
        selectFirst("subgroups")
        //loadSubgroupPupils();
    })
    .done(function() {
        
    })
    .fail(function() {
        showMessage('Fail in loading pupils');
    })
    .always(function() {
    });
}

function loadSubgroupPupils() {
    var sel = getSelectedOption("subgroups")
    emptyList("groupsPupils");
    if (sel != undefined) {
        $.get("/api/subgroup/pupils?task="+taskID + "&groupId="+sel.val() , function(json) {
            
            $.each(JSON.parse(json), function(i, value) {
                var name = $("#pupils option[value='"+ value.refId + "']").text()
                $('#groupsPupils').append($('<option>').text(name).attr('value', value.refId)) ;
            });

            $("#isUnite").prop('checked', sel.data('isUnite'));
            $("#isGenderSensitive").prop('checked', sel.data('isGenderSensitive'));
            $("#isSpreadEvenly").prop('checked', sel.data('isSpreadEvenly'));
            $("#isInactive").prop('checked', sel.data('isInactive'));
            $("#maxAllowed").val(sel.data('maxAllowed'));
            $("#minAllowed").val(sel.data('minAllowed'));
        })
        .done(function() {
            
        })
        .fail(function() {
            showMessage('Fail in loading pupils');
        })
        .always(function() {
        });
    } 
}

function loadTasks() {
    $.get("/api/tasks" , function(json) {
        
        emptyList("tasks")
        $.each(JSON.parse(json), function(i, value) {
            $('#tasks').append($('<option>').text(value.name).attr('value', value.id));
        });
            
        selectFirst("tasks")
    })
    .done(function() {
        
    })
    .fail(function() {
        showMessage('Fail in loading tasks');
    })
    .always(function() {
    });
}

function loadAvailableResults() {
    $.get("/api/available-results?task="+taskID , function(json) {
        
        emptyList("results")
        $.each(JSON.parse(json), function(i, value) {
            $('#results').append($('<option>').text(value.runDate + value.title).attr('value', value.id));
        });
            
        selectFirst("results")
    })
    .done(function() {
        
    })
    .fail(function() {
        showMessage('Fail in loading tasks');
    })
    .always(function() {
    });
}

function saveGroupsPupils() {
    
    var sel = getSelectedOption("subgroups")
    if (sel != undefined) {
        var groupId = sel.val()
        var pupils = []
        $("#groupsPupils option").each(function(i) {
            var opt = $("#groupsPupils option").eq(i);
            pupils.push({"id":groupId, "refId":opt.val()})
        });
        var jsonStr = JSON.stringify(pupils);
        
        $.post("/api/subgroup/pupils?task="+taskID + "&groupId="+groupId, jsonStr, function(json) {
            showMessage('Successfully saved');

        })
        .done(function() {
            
        })
        .fail(function() {
            showMessage('Fail in loading pupils');
        })
        .always(function() {
        });

    }
}