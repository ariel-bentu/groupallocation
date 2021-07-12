
import React, { useState } from 'react';
import { List, ListItem, ListItemText, TextField, SvgIcon } from '@material-ui/core';
import useStyles from "./styles.js"
import { VBox, Spacer, Header } from './elems'

export default function SearchList(props) {
    const classes = useStyles();
    const [search, setSearch] = useState([]);

    let actualItems = props.items ? props.items.filter(i => i.name.includes(search)) : [];

    return (
        <VBox>
            <Header>{props.title}</Header>
            <TextField id="standard-search" type="search"
                helperText="הכנס חלק של שם לחיפוש"
                onChange={(e) => setSearch(e.currentTarget.value)}>{search}</TextField>
            <Spacer />
            <List className={classes.list} style={{ height: 250 }}>
                {actualItems.map((item) => (
                    <ListItem className={classes.listItem} key={item.id}
                        button selected={props.current === item.id}
                        onClick={() => props.onSelect(item.id)}
                        onDoubleClick={() => props.onDoubleClick(item.id)}
                    >

                        {props.genderIcon?<SvgIcon>
                            {item.isMale ?
                                <path d="M17.5,9.5C17.5,6.46,15.04,4,12,4S6.5,6.46,6.5,9.5c0,2.7,1.94,4.93,4.5,5.4V17H9v2h2v2h2v-2h2v-2h-2v-2.1 C15.56,14.43,17.5,12.2,17.5,9.5z M8.5,9.5C8.5,7.57,10.07,6,12,6s3.5,1.57,3.5,3.5S13.93,13,12,13S8.5,11.43,8.5,9.5z" /> :
                                <path d="M9.5,11c1.93,0,3.5,1.57,3.5,3.5S11.43,18,9.5,18S6,16.43,6,14.5S7.57,11,9.5,11z M9.5,9C6.46,9,4,11.46,4,14.5 S6.46,20,9.5,20s5.5-2.46,5.5-5.5c0-1.16-0.36-2.23-0.97-3.12L18,7.42V10h2V4h-6v2h2.58l-3.97,3.97C11.73,9.36,10.66,9,9.5,9z" />
                            }
                        </SvgIcon>:null}

                        <ListItemText className={classes.listItemText} primary={item.name} />
                    </ListItem>
                ))}
            </List>
        </VBox>
    );
}