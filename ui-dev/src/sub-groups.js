import React, { useEffect, useState } from 'react';
import {   Button,  Table, TableBody, TableHead, TableRow, TableCell } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, Paper1, Paper2 } from './elems'
import SearchList from './list-with-search'
import EditGroup from './edit-group'

export default function SubGroups(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);
    const [groups, setGroups] = useState([]);
    const [editGroup, setEditGroup] = useState(undefined);
    const [members, setMembers] = useState(undefined);
    const [editMembers, setEditMembers] = useState(undefined);
    const [editGroupDialog, setEditGroupDialog] = useState(undefined);



    useEffect(() => {
        props.setDirty(editGroup !== undefined || editMembers !== undefined);
    }, [editGroup]); // eslint-disable-line react-hooks/exhaustive-deps

    useEffect(() => {
        api.loadSubGroups(props.currentTask).then(grps => setGroups(grps));
    }, [props.currentTask]);


    const selectSubgroup = (id) => {
        if (dirty()) {
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
            setCurrent(groups.find(p => p.id === id))
        }
    }

    const actGroup = () => editGroup ? editGroup : current;
    const actMembers = () => editMembers ? editMembers : members || [];
    const dirty = () => editMembers || editGroup;

    const addPupilToSubGroup = (pupilId) => {
        if (!current)
            return
        
        if (actMembers().find(p=>p.refId === pupilId)) {
            props.msg.notify("תלמיד זה כבר חבר בקבוצה זו")
            return
        }
        let srcMembers = actMembers();

        let newMembers = [...srcMembers, {
            id: current.id,
            refId: pupilId
        }]
        setEditMembers(newMembers);
    }

    const removePupilFromSubGroup = (pupilId) => {
        let newMembers = actMembers().filter(m => m.refId !== pupilId);
        setEditMembers(newMembers);
    }


    const saveSubGroup = () => {
        api.saveSubGroup(props.currentTask, current.id, editGroup).then(() => {
            props.msg.notify("נשמר בהצלחה");
            api.loadSubGroups(props.currentTask).then(grps => {
                setGroups(grps);
                setEditGroup(undefined);
            })
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
            api.loadSubGroupMembers(props.currentTask, current.id).then(gMembers => {
                setEditMembers(undefined);
                setMembers(gMembers);
            });
        }
    }, [current]); // eslint-disable-line react-hooks/exhaustive-deps


    return (
        <div className={classes.paperContainer}>
            <Paper1>
                <Header>קבוצות{" (" + (groups ? groups.length : "-") + ")"}</Header>
                <VBox>
                    <SearchList items={groups} current={current ? current.id : undefined}
                        style={{height:'60vh', width:'90%'}}
                        onSelect={(id) => selectSubgroup(id)}
                        onDoubleClick={() => { }}
                    />
                    <Spacer/>
                    <HBox>
                    <Button
                        variant="outlined"
                        color="primary"
                        onClick={() => setEditGroupDialog(actGroup())}
                        disabled={!current}
                    >ערוך פרטי קבוצה...</Button>
                    <Spacer/>
                     <Button
                        variant="outlined"
                        color="primary"
                        onClick={() => setEditGroupDialog({})}
                    >קבוצה חדשה...</Button>
                    </HBox>
                </VBox>
            </Paper1>
            {current ?
                <Paper2>
                    <Header>{current.name}{" (" + actMembers().length + ")"}{editMembers ? " - בעריכה" : ""}</Header>
                    <Spacer />
                    <HBox>
                        <VBox>
                            <SearchList items={props.pupils} genderIcon={true}
                                onSelect={() => { }}
                                onDoubleClick={(id) => {
                                    addPupilToSubGroup(id)
                                }}
                                instruction={"דאבל-קליק להוספת תלמיד/ה לקבוצה"}
                            />
                            
                        </VBox>
                        <VBox>
                            <VBox style={{ maxHeight: 400, overflow: 'auto' }}>
                                <Table style={{ width: '100%' }}>
                                    <TableHead>
                                        <TableRow>
                                            <TableCell align="right">שם</TableCell>
                                            <TableCell align="right">פעולות</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {actMembers() ? actMembers().map((p) => (

                                            <TableRow>
                                                <TableCell align="right">{(props.pupils.find(pupil => pupil.id === p.refId) || { name: "" }).name}</TableCell>
                                                <TableCell >
                                                    <HBox>
                                                        <Button variant="outlined" color="primary"
                                                            onClick={() => removePupilFromSubGroup(p.refId)}
                                                        >מחק</Button>
                                                    </HBox>
                                                </TableCell>
                                            </TableRow>

                                        )) : null}
                                    </TableBody>
                                </Table>

                            </VBox>
                            <Spacer />
                            <HBox>
                            
                                <Button variant="outlined" color="primary" disabled={!dirty()}
                                    onClick={() => {
                                        if (editGroup)
                                            saveSubGroup();
                                        if (editMembers)
                                            saveSubGroupMembers();
                                    }}>שמור</Button>
                                <Spacer />
                                {dirty() ? <Button variant="outlined" color="primary"
                                    onClick={() => {
                                        setEditMembers(undefined)
                                        setEditGroup(undefined)
                                    }}>בטל</Button> : null}
                            </HBox>
                        </VBox>


                    </HBox>
                </Paper2> : null}
            <EditGroup open={editGroupDialog !== undefined} group={editGroupDialog}
                Save={newGroup => {
                    setEditGroupDialog(undefined)
                    alert("dodo")
                }}
                Cancel={() => setEditGroupDialog(undefined)}
            />
        </div>

    );
}