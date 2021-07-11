import React from 'react';
import { Box, TextField, Paper } from '@material-ui/core';
import MyAlert from '@material-ui/lab/Alert';



export function VBox (props) {
    return <Box style={{display:'flex', flexDirection:'column', alignItems:'center', ...props.style}}>
        {props.children}
    </Box>
}

export function HBox (props) {
    return <Box style={{display:'flex', flexWrap:'wrap', flexDirection:'row', alignItems:'flex-start', ...props.style}}>
        {props.children}
    </Box>
}

export function Spacer (props) {
    return <dir style={{width:props.width?props.width:5}}/>
}

export function WPaper (props) {
    return <Paper style={{width:'60%'}}>{props.children}</Paper>
}

export function Header (props) {
    return <VBox><dir style={{margin:0, marginBottom:10, height: 25, fontSize:28}}>{props.children}</dir></VBox>;
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

export function Alert(props) {
    return <MyAlert elevation={6} variant="filled" {...props} />;
  }
  

