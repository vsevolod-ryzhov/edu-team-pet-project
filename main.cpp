#include <iostream>
#include <queue>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

bool checking(vector<string> &matrix, int &H, pair<int, int> &START, pair<int, int> &END, pair<int, int> check = {-1, -1}) {
    if (matrix.size() != H) {
        cout << "Error size" << endl;
        return false;
    }
    for (int i = 0; i < H; i++) {
        if (matrix[i].find('S') != -1) {
            if (START != check) {
                START = {i + 1, matrix[i].find('S') + 1};
            }
            else {
                cout << "The most of Start points" << endl;
                return false;
            }
        }
        if (matrix[i].find('G') != -1) {
            if (END != check) {
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
        if (q.front().first - 2 >= 0 && matrix[q.front().first - 2][q.front().second - 1] != '#' && !result.contains({q.front().first - 2, q.front().second - 1})) {
            result[{q.front().first - 1, q.front().second - 1}].push_back({q.front().first - 2, q.front().second - 1});
            q.push({q.front().first - 2, q.front().second - 1});
        }
        if (q.front().second - 2 >= 0 && matrix[q.front().first - 1][q.front().second - 2] != '#' && !result.contains({q.front().first - 1, q.front().second - 2})) {
            result[{q.front().first - 1, q.front().second - 1}].push_back({q.front().first - 1, q.front().second - 2});
            q.push({q.front().first - 1, q.front().second - 2});
        }
        if (q.front().first < H && matrix[q.front().first][q.front().second - 1] != '#' && !result.contains({q.front().first, q.front().second - 1})) {
            result[{q.front().first - 1, q.front().second - 1}].push_back({q.front().first, q.front().second - 1});
            q.push({q.front().first, q.front().second - 1});
        }
        if (q.front().second < W && matrix[q.front().first - 1][q.front().second] != '#' && !result.contains({q.front().first - 1, q.front().second})) {
            result[{q.front().first - 1, q.front().second - 1}].push_back({q.front().first - 1, q.front().second});
            q.push({q.front().first - 1, q.front().second});
        }
        q.pop();
    }
    if (q.empty()) {
        return false;
    }
    return true;
}

int main() {
    int W, H;
    string temp;
    vector<string> matrix;
    cin >> W >> H;
    while(cin >> temp) {
        if (temp.size() < W || temp.size() > W) {
            cout << "Error size" << endl;
            return 0;
        }
        matrix.push_back(temp);
    }
    pair<int, int> check = {-1, -1};
    pair<int, int> START = {-1, -1};
    pair<int, int> END = {-1, -1};
    if (checking(matrix, W, START, END)) {
        map<pair<int, int>, vector<pair<int, int>>> result;
        if (FromMatrixToMap(matrix, START, result, W, H)) {

        }
    }


    return 0;
}