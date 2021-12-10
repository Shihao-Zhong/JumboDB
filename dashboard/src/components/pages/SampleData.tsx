import Box from "@mui/material/Box";
import * as React from "react";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {useDispatch, useSelector} from "react-redux";
import axios from "axios";
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import * as actionTypes from "../../store/actionTypes";


export default function SampleData() {
    const [putValue, setPutValue] = React.useState("");
    const [putKey, setPutKey] = React.useState("");
    const [getKey, setGetKey] = React.useState("");
    const dispatch = useDispatch();

    const connectionState = useSelector((state: any) => state.connection);
    const sampleDataState = useSelector((state: any) => state.sampleData);

    const putData = () => {
        axios.post(`http://${connectionState.allHosts[connectionState.currentIndex]}/resources`, {
            "key": putKey,
            "value": putValue
        }).then((response) => {
            alert("successfully add data into db");
            getAllData();
        });
    }

    const delData = (key: string) => {
        axios.delete(`http://${connectionState.allHosts[connectionState.currentIndex]}/resources/${key}`
        ).then((response) => {
            alert(`successfully delete ${key} data into db`);
            getAllData();
        });
    }

    const getAllData = () => {
        axios.get(`http://${connectionState.allHosts[connectionState.currentIndex]}/resources`).then((response: any) => {
            dispatch({
                type: actionTypes.UPDATE_ALL_DATA,
                payload: response.data.data
            });
        });
    }

    const getData = () => {
        axios.get(`http://${connectionState.allHosts[connectionState.currentIndex]}/resources/${getKey}`).then((response: any) => {
            alert(`successfully query ${getKey} with data ${response.data.value || "NULL"}`)
        });
    }

    React.useEffect(() => {
        getAllData();
    }, []);

    return <>
        <Box component="main" sx={{marginTop: 5, p: 1,}}>

            <Box sx={{p: 1}}>
                <Stack direction="row" spacing={2}>
                    <TextField
                        fullWidth
                        id="outlined-name"
                        label="key"
                        value={putKey}
                        onChange={(event: any) => setPutKey(event.target.value)}
                    />
                    <TextField
                        fullWidth
                        id="outlined-name"
                        label="value"
                        value={putValue}
                        onChange={(event: any) => setPutValue(event.target.value)}
                    />
                    <Button variant="outlined" onClick={putData}>Add</Button>
                </Stack>
            </Box>
            <Box sx={{p: 1}}>
                <Stack direction="row" spacing={2}>
                    <TextField
                        fullWidth
                        id="outlined-name"
                        label="key"
                        value={getKey}
                        onChange={(event: any) => setGetKey(event.target.value)}
                    />

                    <Button variant="outlined" onClick={getData}>Query</Button>
                </Stack>
            </Box>
            <Box sx={{p: 1}}>
                <TableContainer component={Paper}>
                    <Table sx={{minWidth: 650}} aria-label="datatable">
                        <TableHead>
                            <TableRow>
                                <TableCell>key</TableCell>
                                <TableCell>value</TableCell>
                                <TableCell>transactionId</TableCell>
                                <TableCell>operation</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {sampleDataState.allData.map((row: any) => (
                                <TableRow
                                    key={row.key}
                                    sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                >
                                    <TableCell component="th" scope="row" key={`${row.key}-key`}>
                                        {row.key}
                                    </TableCell>
                                    <TableCell align="center" key={row.value}>{row.value}</TableCell>
                                    <TableCell align="center" key={row.transactionId}>{row.transactionId}</TableCell>
                                    <TableCell>
                                        <Button variant="outlined" onClick={() => {delData(row.value)}}>Del</Button>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Box>
        </Box>
    </>
}
