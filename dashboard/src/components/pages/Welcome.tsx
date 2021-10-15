import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import * as React from "react";


export default function Welcome() {

    return (<Box component="main" sx={{flexGrow: 1, p: 3}}>
        <Toolbar/>
        <Typography paragraph>
            Welcome to JumboDB.
        </Typography>
    </Box>)
}
