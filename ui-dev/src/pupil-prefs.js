import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, Button, FormControlLabel, Checkbox, Table, TableBody, TableHead, TableRow, TableCell } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, WPaper } from './elems'
import SearchList from './list-with-search'

export default function PupilPref(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);
    const [prefs, setPrefs] = useState([]);
    const [editPrefs, setEditPrefs] = useState(undefined);

    useEffect(() => {
        props.setDirty(editPrefs !== undefined);
    }, [editPrefs]);


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
        if (srcPrefs.length == 3) {
            props.msg.notify("אין אפשרות להוסיף עוד העדפות")
            return;
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
    }, [current]);

    let visPref = actPrefs();

    return (
        <div className={classes.paperContainer}>
            <Paper elevation={3} className={classes.paper}>
                <Header>תלמידים</Header>
                <VBox>
                    <SearchList items={props.pupils} current={current ? current.id : undefined}

                        onSelect={(id) => selectPupil(id)}
                        onDoubleClick={() => alert("double click")}
                    />
                </VBox>
            </Paper>
            {current ?
                <WPaper>
                    <Header>העדפות של {current.name}{editPrefs ? " - בעריכה" : ""}</Header>
                    <Spacer />
                    <HBox>
                        <SearchList items={props.pupils.filter(p => p.id !== current.id)}
                            onSelect={() => { }}
                            onDoubleClick={(id) => {
                                addPref(id)
                            }}
                        />
                        <VBox>
                            <Table style={{ width: '50%' }} aria-label="simple table">
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
                </WPaper> : null}
        </div>

    );
}