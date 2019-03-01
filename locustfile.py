from locust import HttpLocust, TaskSet, task
import uuid, json
myheaders = {'Content-Type': 'application/json', 'Accept': 'application/json'}

class UserBehavior(TaskSet):
    
    def on_start(self):
        """ on_start is called when a Locust start before any task is scheduled """
        pass

    def on_stop(self):
        """ on_stop is called when the TaskSet is stopping """
        pass

    @task(1)
    def index(self):
        
        self.client.post("/auction", data=json.dumps({"adPlacementId": str(uuid.uuid4())}), headers=myheaders)

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 0
    max_wait = 200