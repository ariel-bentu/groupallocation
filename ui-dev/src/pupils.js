import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, Button, FormControlLabel, Checkbox } from '@material-ui/core';
import useStyles from "./styles.js"
import * as api from './api'
import { VBox, HBox, Spacer, Header, ROField } from './elems'
import SearchList from './list-with-search'

export default function Pupils(props) {
    const classes = useStyles();
    const [current, setCurrent] = useState(undefined);

    return (
        <div className={classes.paperContainer}>
            <Paper elevation={3} className={classes.paper}>
                <Header>תלמידים</Header>
                <VBox>
                    <SearchList items={props.pupils} current={current ? current.id : undefined}

                        onSelect={(id) => setCurrent(props.pupils.find(p => p.id === id))}
                        onDoubleClick={() => alert("double click")}
                    />
                    <Button variant="outlined" color="primary" onClick={() => { }}>הוסף תלמיד...</Button>

                </VBox>
            </Paper>
            {current ? <Paper>
                <Header>פרטים של {current.name}</Header>
                <Spacer />
                <VBox>

                    <ROField label={"שם"} value={current.name} />
                    <ROField label={"מין"} value={current.isMale ? "בן" : "בת"} />
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={current.active}

                                name="checkedB"
                                color="primary" />
                        }
                        label={"פעיל"}
                    />
                    <ROField label={"הערות"} value={current.remarks} />
                    <Spacer />
                    <HBox>
                        <Button variant="outlined" color="primary" onClick={() => { }}>מחק תלמיד</Button>
                        <Button variant="outlined" color="primary" onClick={() => { }}>ערוך תלמיד...</Button>
                    </HBox>
                </VBox>
            </Paper> : null}
        </div>

    );
}