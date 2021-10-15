import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import Divider from "@mui/material/Divider";
import Drawer from "@mui/material/Drawer";
import IconButton from "@mui/material/IconButton";
import CastConnectedIcon from "@mui/icons-material/CastConnected";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import StorageIcon from '@mui/icons-material/Storage';
import MenuIcon from "@mui/icons-material/Menu";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";

import {Link} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import * as actionTypes from "../../store/actionTypes";

const drawerWidth = 240;


export default function LeftOperationBar(props: any) {
    const [mobileOpen, setMobileOpen] = React.useState(false);
    const header = useSelector((state: any) => state.display.header);
    const dispatch = useDispatch();

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
    };

    const updateHeader = (newHeader: string) => {

        dispatch({
            type: actionTypes.UPDATE_UI_HEADER,
            header: newHeader,
        })
    }

    const drawer = (
        <div>
            <Toolbar>
                Jumbo Dashboard
            </Toolbar>
            <Divider/>
            <List>
                <ListItem key="Config Management">
                    <ListItemText primary="Config Management"/>
                </ListItem>
                <ListItem button component={Link} to="/connection" onClick={() => updateHeader("Connection")}>
                    <ListItemIcon>
                        <CastConnectedIcon/>
                    </ListItemIcon>
                    <ListItemText primary="Connection"/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem key="Data Management">
                    <ListItemText primary="Data Management"/>
                </ListItem>
                <ListItem button component={Link} to="/sampleData"
                          onClick={() => updateHeader("Sample data management")}>
                    <ListItemIcon>
                        <StorageIcon/>
                    </ListItemIcon>
                    <ListItemText primary="Sample data"/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem key="Cluster Management">
                    <ListItemText primary="Cluster Management"/>
                </ListItem>
            </List>
        </div>
    );

    return (
        <Box sx={{display: "flex"}}>
            <CssBaseline/>
            <AppBar
                position="fixed"
                sx={{
                    width: {sm: `calc(100% - ${drawerWidth}px)`},
                    ml: {sm: `${drawerWidth}px`},
                }}
            >
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={handleDrawerToggle}
                        sx={{mr: 2, display: {sm: "none"}}}
                    >
                        <MenuIcon/>
                    </IconButton>
                    <Typography variant="h6" noWrap component="div">
                        {header}
                    </Typography>
                </Toolbar>
            </AppBar>
            <Box
                component="nav"
                sx={{width: {sm: drawerWidth}, flexShrink: {sm: 0}}}
                aria-label="mailbox folders"
            >
                {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
                <Drawer
                    variant="temporary"
                    open={mobileOpen}
                    onClose={handleDrawerToggle}
                    ModalProps={{
                        keepMounted: true, // Better open performance on mobile.
                    }}
                    sx={{
                        display: {xs: "block", sm: "none"},
                        "& .MuiDrawer-paper": {boxSizing: "border-box", width: drawerWidth},
                    }}
                >
                    {drawer}
                </Drawer>
                <Drawer
                    variant="permanent"
                    sx={{
                        display: {xs: "none", sm: "block"},
                        "& .MuiDrawer-paper": {boxSizing: "border-box", width: drawerWidth},
                    }}
                    open
                >
                    {drawer}
                </Drawer>
            </Box>
            <Box component="main" sx={{flexGrow: 1, p: 3}}>
                {props.children}
            </Box>
        </Box>
    );
}