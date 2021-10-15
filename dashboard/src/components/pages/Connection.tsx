import * as React from "react";
import Box from "@mui/material/Box";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Divider from "@mui/material/Divider";
import InboxIcon from "@mui/icons-material/Inbox";
import SettingsEthernetIcon from "@mui/icons-material/SettingsEthernet";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import axios from "axios";

import {useDispatch, useSelector} from "react-redux";
import * as actionTypes from "../../store/actionTypes";

export default function Connection() {
    const connectionState = useSelector((state: any) => state.connection);
    const dispatch = useDispatch();
    const currentHost = connectionState.currentIndex !== -1 ?
        connectionState.allHosts[connectionState.currentIndex] : "No host";

    const handleListItemClick = (
        event: React.MouseEvent<HTMLDivElement, MouseEvent>,
        index: number,
    ) => {
        dispatch({
            type: actionTypes.UPDATE_HOST,
            payload: index,
        })
        console.log(connectionState)
    };
    const listItems = connectionState.allHosts.map((host: any, index: any) => {
        return (<ListItemButton key={index}
                                selected={connectionState.currentIndex === index}
                                onClick={(event) => handleListItemClick(event, index)}
        >
            <ListItemText primary={host}/>
        </ListItemButton>)
    });

    const [host, setHost] = React.useState("");
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setHost(event.target.value);
    };

    const addHost = () => {
        axios.get(`http://${host}/health`).then((response: any) => {
            dispatch({
                type: actionTypes.ADD_HOST,
                payload: host
            });
        }).catch((error: any) => {
            console.log(error);
            alert("Connection failed");
        })
    }

    return (
        <>
            <Box component="main" sx={{marginTop: 5, p: 1, display: "flex",}}>
                <Box sx={{width: "100%", maxWidth: 360, bgcolor: "background.paper"}}>
                    <List component="nav" aria-label="main mailbox folders">
                        <ListItem>
                            <ListItemIcon>
                                <InboxIcon/>
                            </ListItemIcon>
                            <ListItemText primary="Connections"/>
                        </ListItem>
                        <ListItem>
                            <ListItemIcon>
                                <SettingsEthernetIcon/>
                            </ListItemIcon>
                            <ListItemText primary={currentHost}/>
                        </ListItem>
                    </List>
                    <Divider/>
                    <List component="nav" aria-label="secondary mailbox folder">
                        {listItems}
                    </List>

                </Box>
                <Box sx={{p: 10}}>
                    <Stack direction="row" spacing={2}>
                        <TextField
                            fullWidth
                            id="outlined-name"
                            label="host:port"
                            value={host}
                            onChange={handleChange}
                        />
                        <Button variant="outlined" onClick={addHost}>Connect</Button>
                    </Stack>

                </Box>

            </Box>
        </>
    );
}