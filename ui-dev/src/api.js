
export async function loadTasks() {
    return fetch( "/api/tasks").then(res=>res.json());
}

export async function loadPupils(task) {
    if (task === "")
        return [];
    return fetch(`/api/pupils?task=${task}`).then(res=>res.json());
}
export async function loadPupilPrefs(task, pupilId) {
    return fetch(`/api/pupil/prefs?task=${task}&pupilId=${pupilId}`).then(res=>res.json());
}

export async function loadSubGroups(task) {
    return fetch(`/api/subgroups?task=${task}`).then(res=>res.json());
}


export async function loadSubGroupMembers(task, groupId) {
    return fetch(`/api/subgroup/pupils?task=${task}&groupId=${groupId}`).then(res=>res.json());
}


export async function loadResults(task) {
    return fetch( `/api/available-results?task=${task}`).then(res=>res.json());
}

export async function saveResultName(task, resId, newName) {
    const requestOptions = {
        method: 'POST',
//        headers: { 'Content-Type': 'application/json' },
//        body: JSON.stringify({ title: 'React PUT Request Example' })
    };
    return fetch("/api/result?task=" + task + "&id=" + resId + "&resultName=" + decodeURIComponent(newName), requestOptions);
}

export async function savePupilPerfs(task, pupilId, prefs) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(prefs)
    };
    return fetch(`/api/pupil/prefs?task=${task}&pupilId=${pupilId}`, requestOptions);
}
/**
 * 
{
    "name": name,
    "isUnite": isUnite,
    "isSpreadEvenly": isSpreadEvenly,
    "isGenderSensitive": isGenderSensitive,
    "isInactive": isInactive,
    "minAllowed": minAllowed,
    "maxAllowed": maxAllowed,
    "isGarden": isGarden
}
 */
export async function saveSubGroup(task, group) {
    const requestOptions = {
        method: group.id === ""?'POST':'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(group)
    };
    return fetch(`/api/subgroups?task=${task}&groupId=${group.id}`, requestOptions);
}
//[{id, refId}]
export async function saveSubGroupMembers(task, groupId, members) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(members)
    };
    return fetch(`/api/subgroup/pupils?task=${task}&groupId=${groupId}`, requestOptions);
}