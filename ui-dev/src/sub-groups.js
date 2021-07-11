import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, Button, FormControlLabel, Checkbox, Table, TableBody, TableHead, TableRow, TableCell } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, WPaper } from './elems'
import SearchList from './list-with-search'

export default function SubGroups(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);
    const [groups, setGroups] = useState([]);
    const [editGroup, setEditGroup] = useState(undefined);
    const [members, setMembers] = useState(undefined);
    const [editMembers, setEditMembers] = useState(undefined);

    useEffect(() => {
        props.setDirty(editGroup !== undefined || editMembers !== undefined);
    }, [editGroup]);

    useEffect(() => {
        api.loadSubGroups(props.currentTask);
    }, [props.currentTask]);

    useEffect(() => {
        if (current) {
            api.loadSubGroupMembers(props.currentTask, current.id);
        }
    }, [current]);

    const selectSubgroup = (id) => {
        if (editGroup || editMembers) {
            props.msg.alert({
                title: "שינויים לא נשמרו",
                message: "לפני החלפת קבוצה יש לשמור או לבטל שינויים",
                buttons: [{
                    label: "שמור",
                    callback: () => {
                        if (editGroup)
                            saveSubGroup()
                        if (editMembers)
                            saveSubGroupMembers()
                        setCurrent(props.groups.find(p => p.id === id))
                    }
                },
                {
                    label: "התעלם משינויים",
                    callback: () => setCurrent(props.groups.find(p => p.id === id))
                },
                {
                    label: "בטל",
                    callback: () => { }
                }]
            })
        } else {
            setCurrent(props.groups.find(p => p.id === id))
        }
    }

    const actGroup = () => editGroup ? editGroup : current;
    const actMembers = () => editMembers ? editMembers : members || [];

    const addPupilToSubGroup = (id) => {
        if (!current)
            return
        // let srcPrefs = actPrefs();
        // 
        // let newPrefs = [...srcPrefs, {
        //     id: current.id,
        //     refId,
        //     active: true,
        //     priority: srcPrefs.length
        // }]
        // setEditPrefs(newPrefs);
    }

    const removePupilFromSubGroup = (index) => {
        // let newPrefs = actPrefs().filter((v, i) => i !== index);
        // setEditPrefs(newPrefs);
    }


    const saveSubGroup = () => {
        api.saveSubGroup(props.currentTask, current.id, editGroup).then(() => {
            props.msg.notify("נשמר בהצלחה");
            setGroup(editGroup);
            setEditGroups(undefined)
        })
    }

    const saveSubGroupMembers = () => {
        api.saveSubGroupMembers(props.currentTask, current.id, editMembers).then(() => {
            props.msg.notify("נשמר בהצלחה");
            setMembers(editMembers);
            setEditMembers(undefined)
        })
    }

    useEffect(() => {
        if (current) {
            console.log("Load members")
            api.loadSubGroupPupils(props.currentTask, current.id).then(gr => {
                setEditGroups(undefined);
                setGroups(gr);
            });
        }
    }, [current]);


    return (
        <div className={classes.paperContainer}>
            <Paper elevation={3} className={classes.paper}>
                <Header>קבוצות</Header>
                <VBox>
                    <SearchList items={groups} current={current ? current.id : undefined}
                        onSelect={(id) => selectSubgroup(id)}
                        onDoubleClick={()=>{}}
                    />
                </VBox>
            </Paper>
            {current ?
                <WPaper>
                    <Header>{current.name}{editMembers ? " - בעריכה" : ""}</Header>
                    <Spacer />
                    <HBox>
                        <SearchList items={props.pupils}
                            onSelect={() => { }}
                            onDoubleClick={(id) => {
                                addPupilToSubGroup(id)
                            }}
                        />
                        <VBox>
                        <SearchList items={actMembers()}
                            onSelect={() => { }}
                            onDoubleClick={(id) => {
                                removePupilFromSubGroup(id)
                            }}
                        />
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