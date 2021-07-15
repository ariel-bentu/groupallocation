import React from 'react';
import { Box, TextField, Paper, FormControlLabel, Checkbox, Button } from '@material-ui/core';
import MyAlert from '@material-ui/lab/Alert';



export function VBox(props) {
    return <Box style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', ...props.style }}>
        {props.children}
    </Box>
}

export function HBox(props) {
    return <Box style={{ display: 'flex', flexWrap: 'wrap', flexDirection: 'row', alignItems: 'flex-start', ...props.style }}>
        {props.children}
    </Box>
}


export function Spacer(props) {
    return <dir style={{ width: props.width ? props.width : 5 }} />
}
export function GButton(props) {
    return <Button variant="outlined" color="primary" style={{margin:3}} {...props}>{props.label}</Button>
}

export function Text(props) {
    return <dir style={{ fontSize: 12 }} >{props.children}</dir>
}


export function Paper1(props) {
    return <Paper elevation={3} style={{ width: props.width ? props.width : '27%', height: '85vh', ...props.style }} {...props}>{props.children}</Paper>
}

export function Paper2(props) {
    return <Paper1 width={props.width || '60%'}>{props.children}</Paper1>
}

export function Header(props) {
    return <VBox><dir style={{ margin: 0, marginBottom: 10, height: 25, fontSize: 28 }}>{props.children}</dir></VBox>;
}

export function ROField(props) {
    return <TextField
        label={props.label}
        value={props.value}
        InputProps={{
            readOnly: true,
        }}
    />
}

export function LabeledCheckBox(props) {
    return <FormControlLabel
        control={
            <Checkbox
                checked={props.checked}
                onClick={props.onClick}
                color="primary"
            />
        }
        labelPlacement="start"
        label={props.label}
    />
}
export function Alert(props) {
    return <MyAlert elevation={6} variant="filled" {...props} />;
}


