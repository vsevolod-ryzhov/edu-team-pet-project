#include <iostream>
#include <queue>
#include <vector>
#include <algorithm>
#include <map>
#include <fstream>

using namespace std;

bool checking(vector<string> &matrix, int &H, pair<int, int> &START, pair<int, int> &END, pair<int, int> check = {-1, -1}) {
    if (matrix.size() != H) {
        cout << "Error size" << endl;
        return false;
    }
    for (int i = 0; i < H; i++) {
        if (matrix[i].find('S') != -1) {
            if (START == check) {
                START = {i + 1, matrix[i].find('S') + 1};
            }
            else {
                cout << "The most of Start points" << endl;
                return false;
            }
        }
        if (matrix[i].find('G') != -1) {
            if (END == check) {
                END = {i + 1, matrix[i].find('G') + 1};
            }
            else {
                cout << "The most of Finish points" << endl;
                return false;
            }
        }
    }
    if (END == check || START == check) {
        cout << "No Finish point or Start point" << endl;
        return false;
    }
    return true;
}

bool FromMatrixToMap(vector<string> &matrix, pair<int, int> START, map<pair<int, int>, vector<pair<int, int>>> &result, int H, int W) {
    queue<pair<int, int>> q {{START}};
    result[START];
    while (!q.empty() && matrix[q.front().first - 1][q.front().second - 1] != 'G') {
        if (q.front().first - 2 >= 0 && matrix[q.front().first - 2][q.front().second - 1] != '#' && !result.count({q.front().first - 1, q.front().second})) {
            result[{q.front().first, q.front().second}].push_back({q.front().first - 1, q.front().second});
            q.push({q.front().first - 1, q.front().second});
        }
        if (q.front().second - 2 >= 0 && matrix[q.front().first - 1][q.front().second - 2] != '#' && !result.count({q.front().first, q.front().second - 1})) {
            result[{q.front().first, q.front().second}].push_back({q.front().first, q.front().second - 1});
            q.push({q.front().first, q.front().second - 1});
        }
        if (q.front().first < H && matrix[q.front().first][q.front().second - 1] != '#' && !result.count({q.front().first + 1, q.front().second})) {
            result[{q.front().first, q.front().second}].push_back({q.front().first + 1, q.front().second});
            q.push({q.front().first + 1, q.front().second});
        }
        if (q.front().second < W && matrix[q.front().first - 1][q.front().second] != '#' && !result.count({q.front().first, q.front().second + 1})) {
            result[{q.front().first, q.front().second}].push_back({q.front().first, q.front().second + 1});
            q.push({q.front().first, q.front().second + 1});
        }
        q.pop();
    }
    if (q.empty()) {
        return false;
    }
    return true;
}

void print(map<pair<int, int>, pair<int, int>> result, pair<int, int> START, pair<int, int> END) {
    pair<int, int> i = {END.first, END.second};
    vector<pair<int, int>> res;
    while (i != START) {
        res.push_back(i);
        i = result[i];
    }
    cout << START.first << " " << START.second << endl;
    for (int j = res.size() - 1; j >= 0; j--) {
        cout << res[j].first << " " << res[j].second << endl;
    }
}

void BreadthFirstSearch(map<pair<int, int>, vector<pair<int, int>>> &graph, pair<int, int> START, pair<int, int> END) {
    map<pair<int, int>, pair<int, int>> result;
    queue<pair<int, int>> q {{START}};
    while (!q.empty()) {
        for (auto i : graph[q.front()]) {
            result[i] = q.front();
            q.push(i);
            if (i == END) {
                print(result, START, END);
                return;
            }
        }
        q.pop();
    }
}

int main() {
    int W, H;
    string temp;
    vector<string> matrix;
    ifstream in("testdata\\corner.txt");
    if (!in.is_open()) {
        std::cerr << "Error file" << std::endl;
        return 1;
    }
    in >> W >> H;
    while(in >> temp) {
        if (temp.size() < W || temp.size() > W) {
            cout << "Error size" << endl;
            return 0;
        }
        matrix.push_back(temp);
    }
    pair<int, int> check = {-1, -1};
    pair<int, int> START = {-1, -1};
    pair<int, int> END = {-1, -1};
    if (checking(matrix, H, START, END)) {
        map<pair<int, int>, vector<pair<int, int>>> result;
        if (FromMatrixToMap(matrix, START, result, H, W)) {
            BreadthFirstSearch(result, START, END);
        }
        else {
            cout << "There`s no way" << endl;
            return 1;
        }
    }


    return 0;
}