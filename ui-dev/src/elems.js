import React from 'react';
import { Box, TextField, Paper, FormControlLabel, Checkbox } from '@material-ui/core';
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

export function VBoxScroll(props) {
    return <Box style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', ...props.style }}>
        {props.children}
    </Box>
}

export function Spacer(props) {
    return <dir style={{ width: props.width ? props.width : 5 }} />
}



export function WPaper(props) {
    return <Paper style={{ width: '60%' }}>{props.children}</Paper>
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


