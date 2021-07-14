
export async function loadTasks() {
    return fetch("/api/task").then(res => res.json());
}

export async function loadPupils(task) {
    if (task === "")
        return [];
    return fetch(`/api/pupils?task=${task}`).then(res => res.json());
}
export async function loadPupilPrefs(task, pupilId) {
    return fetch(`/api/pupil/prefs?task=${task}&pupilId=${pupilId}`).then(res => res.json());
}

export async function loadSubGroups(task) {
    return fetch(`/api/subgroup?task=${task}`).then(res => res.json());
}


export async function loadSubGroupMembers(task, groupId) {
    return fetch(`/api/subgroup/pupils?task=${task}&groupId=${groupId}`).then(res => res.json());
}

export async function loadPupilSubgroups(task, pupilId) {
    return fetch(`/api/subgroups/pupil?task=${task}&pupilId=${pupilId}`).then(res => res.json());
}        

export async function loadResults(task) {
    return fetch(`/api/result?task=${task}`).then(res => res.json());
}

export async function saveTask( task) {
    const requestOptions = {
        method: task.id === "" ? 'POST' : 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(task)
    };
    return fetch(`/api/task`, requestOptions);
}

export async function saveResultName(task, resId, newName) {
    const requestOptions = {
        method: 'POST',
        //        headers: { 'Content-Type': 'application/json' },
        //        body: JSON.stringify({ title: 'React PUT Request Example' })
    };
    return fetch("/api/result?task=" + task + "&id=" + resId + "&resultName=" + decodeURIComponent(newName), requestOptions);
}

export async function savePupil(task, pupil) {
    console.log(JSON.stringify(pupil))
    const requestOptions = {
        method: task.id === "" ? 'POST' : 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(pupil)
     };
    return fetch(`/api/pupils?task=${task}`, requestOptions);
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
        method: group.id === "" ? 'POST' : 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(group)
    };
    return fetch(`/api/subgroup?task=${task}&groupId=${group.id}`, requestOptions);
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

export function showResults(task, resultId, clean) {
    let cleanStr = clean ? "&clean=true" : "";
    window.open(`/results?task=${task}&id=${resultId}${cleanStr}`);
}

export function deleteTask(task) {
    const requestOptions = {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' }
    };
    return fetch(`/api/task?task=${task}`, requestOptions);
}

export function deleteResult(resultId) {
    const requestOptions = {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' }
    };
    return fetch(`/api/result?id=${resultId}`, requestOptions);
}

export function duplicateResult(task, resultId, newName) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
    };
    return fetch(`/api/result/duplicate?task=${task}&id=${resultId}&resultName=${newName}`, requestOptions);	
}

export function deletePupil(task, pupilId) {
    const requestOptions = {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' }
    };
    return fetch(`/api/pupils?task=${task}&id=${pupilId}`, requestOptions);
}

export function deleteSubgroup(task, groupId) {
    const requestOptions = {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' }
    };
    return fetch(`/api/subgroup?task=${task}&groupId=${groupId}`, requestOptions);
}