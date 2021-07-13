import React, { useEffect, useState } from 'react';
import {
    FormControlLabel, Checkbox, Table, TableBody, TableHead, TableRow, TableCell,
    List, ListItem, ListItemText
} from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, ROField, Paper1, Paper2, GButton } from './elems'
import SearchList from './list-with-search'
import EditPupil from './edit-pupil'


export default function PupilPref(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);
    const [prefs, setPrefs] = useState([]);
    const [editPrefs, setEditPrefs] = useState(undefined);
    const [editPupilDialog, setEditPupilDialog] = useState(undefined);
    const [pupilSubgroups, setPupilSubgroups] = useState([]);


    useEffect(() => {
        props.setDirty(editPrefs !== undefined);
    }, [editPrefs]); // eslint-disable-line react-hooks/exhaustive-deps


    const selectPupil = (id) => {
        if (editPrefs) {
            props.msg.alert({
                title: "שינויים לא נשמרו",
                message: "לפני החלפת התלמיד הנוכחי יש לשמור או לבטל שינויים",
                buttons: [{
                    label: "שמור",
                    callback: () => {
                        save()
                        setCurrent(props.pupils.find(p => p.id === id))
                    }
                },
                {
                    label: "התעלם משינויים",
                    callback: () => setCurrent(props.pupils.find(p => p.id === id))
                },
                {
                    label: "בטל",
                    callback: () => { }
                }]
            })
        } else {
            setCurrent(props.pupils.find(p => p.id === id))
        }
    }

    const actPrefs = () => editPrefs ? editPrefs : prefs;
    const swapPref = (src, target) => {
        let newPrefs = [];
        let srcPrefs = actPrefs();

        for (let i = 0; i < srcPrefs.length; i++) {
            if (i === src) {
                newPrefs[i] = srcPrefs[target];
            } else if (i === target) {
                newPrefs[i] = srcPrefs[src];
            } else {
                newPrefs[i] = srcPrefs[i]
            }
            newPrefs[i].priority = i
        }
        setEditPrefs(newPrefs);
    }

    const addPref = (refId) => {
        let srcPrefs = actPrefs() || [];
        if (srcPrefs)
            if (srcPrefs.length === 3) {
                props.msg.notify("אין אפשרות להוסיף עוד העדפות")
                return;
            }

        if (srcPrefs.find(p => p.refId === refId)) {
            props.msg.notify("תלמיד זה כבר נמצא כהעדפה")
            return
        }

        let newPrefs = [...srcPrefs, {
            id: current.id,
            refId,
            active: true,
            priority: srcPrefs.length
        }]
        setEditPrefs(newPrefs);
    }

    const removePref = (index) => {
        let newPrefs = actPrefs().filter((v, i) => i !== index);
        setEditPrefs(newPrefs);
    }

    const toggleActive = (refId) => {
        let newPrefs = actPrefs().map(r => r.refId === refId ? { id: r.id, refId: r.refId, active: !r.active, priority: r.priority } : r);
        setEditPrefs(newPrefs);
    }

    const save = () => {
        api.savePupilPerfs(props.currentTask, current.id, editPrefs).then(() => {
            props.msg.notify("נשמר בהצלחה");
            setPrefs(editPrefs);
            setEditPrefs(undefined)
        })
    }

    useEffect(() => {
        if (current) {
            console.log("Load preferences")
            api.loadPupilPrefs(props.currentTask, current.id).then(ps => {
                setEditPrefs(undefined);
                setPrefs(ps);
            });

            console.log("Load Groups");
            api.loadPupilSubgroups(props.currentTask, current.id).then((gprs => setPupilSubgroups(gprs)));
        }
    }, [current]); // eslint-disable-line react-hooks/exhaustive-deps

    let visPref = actPrefs();

    return (
        <div className={classes.paperContainer}>
            <Paper1 width='25%'>
                <Header>תלמידים</Header>

                <SearchList items={props.pupils} current={current ? current.id : undefined} genderIcon={true}
                    style={{ width: '80%', height: '85%' }}

                    onSelect={(id) => selectPupil(id)}
                    onDoubleClick={() => { }}
                />
                <GButton label="הוסף תלמיד..." onClick={() => setEditPupilDialog({ active: true })} />

            </Paper1>
            {current ?

                <Paper1 width='25%'>
                    <Header>פרטים של {current.name}</Header>
                    <Spacer />
                    <VBox>

                        <ROField label={"שם"} value={current.name} />
                        <ROField label={"מין"} value={current.isMale ? "בן" : "בת"} />
                        <FormControlLabel
                            control={
                                <Checkbox
                                    checked={current.active}
                                    color="primary" />
                            }
                            label={"פעיל"}
                        />
                        <ROField label={"הערות"} value={current.remarks} />
                        <Spacer />
                        <HBox>
                            <GButton label="מחק תלמיד"
                                onClick={() =>
                                    props.msg.alert({
                                        title: "מחיקת תלמיד",
                                        message: `האם למחוק את התלמיד ${current.name} \nמחיקת התלמיד הינה בלתי הפיכה!!!`,
                                        buttons: [{
                                            label: "מחק",
                                            callback: () => {
                                                api.deletePupil(props.currentTask, current.id).then(() => {
                                                    props.reloadPupil();
                                                    props.msg.notify("תלמיד נמחק בהצלחה")
                                                })
                                            }
                                        },
                                        {
                                            label: "בטל",
                                            callback: () => { }
                                        }]
                                    }
                                    )} disabled={!current} />
                            <GButton label="ערוך תלמיד..." onClick={() => setEditPupilDialog(current)} />
                        </HBox>
                    </VBox>
                </Paper1> : null}

            {current ? <Paper2 width='45%' >
                <Header>העדפות של {current.name}{editPrefs ? " - בעריכה" : ""}</Header>
                <Spacer />
                <HBox  >
                    <SearchList items={props.pupils.filter(p => p.id !== current.id)} genderIcon={true}
                        style={{ width: '40%' }}
                        onSelect={() => { }}
                        onDoubleClick={(id) => {
                            addPref(id)
                        }}
                        instruction={"דאבל-קליק להוספת העדפה"}
                    />
                    <VBox style={{ width: '60%' }}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>#</TableCell>
                                    <TableCell align="right">שם</TableCell>
                                    <TableCell align="right">פעיל</TableCell>
                                    <TableCell align="right">פעולות</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {visPref ? visPref.map((p, index) => (

                                    <TableRow key={index}>
                                        <TableCell component="th" scope="row">{index + 1}</TableCell>
                                        <TableCell align="right">{props.pupils.find(pupil => pupil.id === p.refId).name}</TableCell>
                                        <TableCell align="right">
                                            <Checkbox
                                                checked={p.active}
                                                color="primary"
                                                onClick={(e) => toggleActive(p.refId)}
                                            />
                                        </TableCell>
                                        <TableCell >
                                            <HBox>
                                                <GButton label="מחק" onClick={() => removePref(index)} />
                                                {index > 0 ? <GButton label="&uarr;" onClick={() => swapPref(index, index - 1)} /> : null}
                                                {prefs && index < prefs.length - 1 ? <GButton label="&darr;" onClick={() => swapPref(index, index + 1)} /> : null}
                                            </HBox>
                                        </TableCell>
                                    </TableRow>

                                )) : null}
                            </TableBody>
                        </Table>
                        <HBox>
                            <GButton label="שמור" disabled={!editPrefs} onClick={() => save()} />
                            <Spacer />
                            {editPrefs ? <GButton label="בטל" onClick={() => setEditPrefs(undefined)} /> : null}
                        </HBox>
                        <Header>קבוצות</Header>
                        <List dense={true}>
                            {pupilSubgroups.map((g, index) => (<ListItem key={index}>
                                <ListItemText
                                    primary={g.name}
                                    secondary={g.isUnite ? "איחוד" : "פירוד"}
                                />
                            </ListItem>))}
                        </List>

                    </VBox>
                </HBox>
            </Paper2>
                : null}
            <EditPupil open={editPupilDialog !== undefined} pupil={editPupilDialog}
                Save={newPupil => {
                    setEditPupilDialog(undefined)
                    api.savePupil(props.currentTask, newPupil).then(() => {
                        props.reloadPupil();
                        props.msg.notify(`${newPupil.name} עודכן בהצלחה`)
                    })
                }}
                Cancel={() => setEditPupilDialog(undefined)}
            />
        </div>

    );
}