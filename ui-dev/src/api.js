
export async function loadTasks() {
    return fetch( "/api/tasks").then(res=>res.json());
}

export async function loadPupils(task) {
    return fetch(`/api/pupils?task=${task}`).then(res=>res.json());
}
export async function loadPupilPrefs(task, pupilId) {
    return fetch(`/api/pupil/prefs?task=${task}&pupilId=${pupilId}`).then(res=>res.json());
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