import express from 'express';
import {PangeaConnect} from './middleware/pangea-connect.mjs';
const app = express();
const port = 8080

app.use(PangeaConnect)
app.get('/get', (req, res) => {
  res.send('Hello World!')
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})