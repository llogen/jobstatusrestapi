const express = require('express');
const LimitingMiddleware = require('limiting-middleware');
const fs = require('fs')
const jobArray = require('./jobstatus.json');
const app = express();
const bodyParser = require('body-parser');

const readJobStatus = (ID) => {
    var data = fs.readFileSync('./jobstatus.json', 'utf8', function (err){
        if (err){
            console.log(err);
            return "could not read the db data"
        }
    });
    
    var jsonData = JSON.parse(data)    
    for (var i = 0; i < jsonData.length; i++){
        if (jsonData[i].ID == ID){
            return jsonData[i].Status;
        }
    }
    return "Job does not exist";
};

const updateJobStatus = (req) => {
    var found = 0

    var data = fs.readFileSync('./jobstatus.json', 'utf8', function (err){
        if (err){
            console.log(err);
            return "could not read the db data"
        }
    });
    var jsonData = JSON.parse(data)

    for (var i = 1; i < jsonData.length+1; i++){
        if (jsonData[i-1].ID == req.ID){
            found = i
        }
    }
    console.log(req)
    if (found > 0){
        jsonData[found-1].Status = req.Status
        fs.writeFile('./jobstatus.json', JSON.stringify(jsonData, undefined, 2), function(err){
            if (err) {
                return "could not update the status"
            }
        });
        return "Status edited"
    }else {
        return "Job does not exist"; 
    }
};

const addJobStatus = (req) => {
    var jobExists = ""
    jobExists = readJobStatus(req.ID)
    if (jobExists != "Job does not exist"){
        return "The Job already exist!"
    }
    var data = fs.readFileSync('./jobstatus.json', 'utf8', function (err){
        if (err){
            console.log(err);
            return "could not read the db data"
        }
    });
    var jsonData = JSON.parse(data)
    var appendData = {ID: req.ID, Status: req.Status}
    jsonData.push(appendData)

    fs.writeFile('./jobstatus.json', JSON.stringify(jsonData, undefined, 2), function(err){
        if (err) return "Job could not be added"
    });
    return "Job added"; 
};

const removeJobStatus = (ID) => {
    var jobExists = ""
    jobExists = readJobStatus(ID)
    if (jobExists == "Job does not exist"){
        return "The Job doesn't even exist!"
    }
    var data = fs.readFileSync('./jobstatus.json', 'utf8', function (err){
        if (err){
            return "could not read the db data"
        }
    });
    var jsonData = JSON.parse(data)
    
    for (var i = 0; i < jsonData.length; i++){
        if (jsonData[i].ID == ID){
            jsonData.splice(i,i);
        }
    }

    fs.writeFile('./jobstatus.json', JSON.stringify(jsonData, undefined, 2), function(err){
        if (err) return "Job could not be removed"
    });
    return "Job removed"
};

app.use(express.json());
app.use(express.urlencoded({extended: true}));
app.use(new LimitingMiddleware().limitByIp());

app.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    next();
});

app.get('/', (req, res) => {
    res.send('Try /readjobstatus');
});

app.get('/ping', (req, res) => {
    res.send('pong');
});

app.get('/readjobstatus/:id', (req, res) => {
    res.json(readJobStatus(req.params.id));
});

app.post('/updatejobstatus/', (req, res) => {
    res.json(updateJobStatus(req.body));
});

app.post('/addjobstatus/', (req, res) => {
    res.json(addJobStatus(req.body));
});

app.post('/removejobstatus/:id', (req, res) => {
    res.json(removeJobStatus(req.params.id));
});

app.use((err, req, res, next) => {
    const statusCode = err.statusCode || 500;

    res.status(statusCode).json({
        type: 'error', message: err.message
    });
});

const PORT = process.env.PORT || 3005;
app.listen(PORT, () => console.log(`listening on ${PORT}`));