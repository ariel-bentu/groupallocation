import React, { useEffect, useState } from 'react';
import {  Button, FormControlLabel, Checkbox, Table, TableBody, TableHead, TableRow, TableCell } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, ROField, Paper1, Paper2 } from './elems'
import SearchList from './list-with-search'
import EditPupil from './edit-pupil'


export default function PupilPref(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);
    const [prefs, setPrefs] = useState([]);
    const [editPrefs, setEditPrefs] = useState(undefined);
    const [editPupilDialog, setEditPupilDialog] = useState(undefined);

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
        let srcPrefs = actPrefs();
        if (srcPrefs.length === 3) {
            props.msg.notify("אין אפשרות להוסיף עוד העדפות")
            return;
        }

        if (srcPrefs.find(p=>p.refId === refId)) {
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
                    onDoubleClick={() =>{}}
                />
                <Button variant="outlined" color="primary" 
                                onClick={() => setEditPupilDialog({})}>הוסף תלמיד...</Button>

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
                            <Button variant="outlined" color="primary" onClick={() => { }}>מחק תלמיד</Button>
                            <Button variant="outlined" color="primary" onClick={() => {
                                setEditPupilDialog(current)
                             }}>ערוך תלמיד...</Button>
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
                    <VBox style={{ width: '60%'}}>
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
                                                <Button variant="outlined" color="primary"
                                                    onClick={() => removePref(index)}
                                                >מחק</Button>
                                                {index > 0 ? <Button variant="outlined" color="primary"
                                                    onClick={() => swapPref(index, index - 1)}
                                                >&uarr;</Button> : null}
                                                {index < prefs.length - 1 ? <Button variant="outlined" color="primary"
                                                    onClick={() => swapPref(index, index + 1)}
                                                >&darr;</Button> : null}
                                            </HBox>
                                        </TableCell>
                                    </TableRow>

                                )) : null}
                            </TableBody>
                        </Table>
                        <HBox>
                            <Button variant="outlined" color="primary" disabled={!editPrefs}
                                onClick={() => save()}>שמור</Button>
                            <Spacer />
                            {editPrefs ? <Button variant="outlined" color="primary"
                                onClick={() => setEditPrefs(undefined)}>בטל</Button> : null}
                        </HBox>
                    </VBox>
                </HBox>
            </Paper2>
                : null}
            <EditPupil open={editPupilDialog !== undefined} pupil={editPupilDialog}
                Save={newPupil => {
                    setEditPupilDialog(undefined)
                    alert("todo")
                }}
                Cancel={() => setEditPupilDialog(undefined)}
            />
        </div>

    );
}